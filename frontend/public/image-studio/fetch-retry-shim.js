// Sub2API image-studio fetch retry shim.
//
// Wraps window.fetch so that requests carrying X-Sub2API-Image-Task-Key are
// transparently retried when the browser/network drops the connection
// (TypeError: Failed to fetch / AbortError after long upstream latency).
//
// The matching server-side cache (handler.handleStudioImageTaskRequest) keeps
// the in-flight upstream call alive via context.WithoutCancel and serves the
// final response on any later request that reuses the same task key, so a
// retry that lands while the original is still running will block until the
// cached completion is ready and replay it.
(function () {
  if (typeof window === 'undefined' || !window.fetch || window.__sub2apiImageStudioFetchShim__) {
    return;
  }
  window.__sub2apiImageStudioFetchShim__ = true;

  var HEADER = 'x-sub2api-image-task-key';
  var MAX_ATTEMPTS = 6;
  var BASE_DELAY_MS = 2000;
  var MAX_DELAY_MS = 15000;
  var original = window.fetch.bind(window);

  function readHeader(init) {
    if (!init || !init.headers) return null;
    var h = init.headers;
    if (typeof Headers !== 'undefined' && h instanceof Headers) {
      return h.get(HEADER);
    }
    if (Array.isArray(h)) {
      for (var i = 0; i < h.length; i++) {
        if (Array.isArray(h[i]) && typeof h[i][0] === 'string' && h[i][0].toLowerCase() === HEADER) {
          return h[i][1];
        }
      }
      return null;
    }
    if (typeof h === 'object') {
      var keys = Object.keys(h);
      for (var j = 0; j < keys.length; j++) {
        if (keys[j].toLowerCase() === HEADER) return h[keys[j]];
      }
    }
    return null;
  }

  function detectTaskKey(input, init) {
    try {
      if (typeof Request !== 'undefined' && input instanceof Request) {
        return input.headers.get(HEADER);
      }
    } catch (_) {}
    return readHeader(init);
  }

  function isRetriable(err) {
    if (!err) return false;
    // Browsers throw TypeError for connection-level failures and AbortError
    // when the timer or upstream cancel was set. Both are safe to retry
    // because the server-side cache deduplicates by task key.
    return err.name === 'TypeError' || err.name === 'AbortError';
  }

  function delayMs(attempt) {
    var d = BASE_DELAY_MS * Math.pow(1.5, attempt - 1);
    return Math.min(d, MAX_DELAY_MS);
  }

  function sleep(ms) {
    return new Promise(function (resolve) { setTimeout(resolve, ms); });
  }

  async function attempt(input, init, taskKey, n) {
    // Request bodies are single-shot streams. Clone before the fetch so a
    // failed attempt can be retried without "body stream already read".
    var nextInput = input;
    var nextInit = init;
    try {
      if (typeof Request !== 'undefined' && input instanceof Request) {
        nextInput = input.clone();
      }
    } catch (_) {}
    try {
      var res = await original(input, init);
      return res;
    } catch (err) {
      if (!isRetriable(err) || n >= MAX_ATTEMPTS) throw err;
      try {
        if (window && window.console && typeof console.warn === 'function') {
          console.warn('[image-studio retry] taskKey=' + taskKey + ' attempt=' + n + ' err=' + err.name + ': ' + err.message);
        }
      } catch (_) {}
      await sleep(delayMs(n));
      return attempt(nextInput, nextInit, taskKey, n + 1);
    }
  }

  window.fetch = function (input, init) {
    var taskKey;
    try {
      taskKey = detectTaskKey(input, init);
    } catch (_) {
      taskKey = null;
    }
    if (!taskKey) {
      return original(input, init);
    }
    return attempt(input, init, taskKey, 1);
  };
})();
