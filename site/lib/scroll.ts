// Offset from the top of the viewport when scrolling a section into view.
// Keeps the section heading comfortably near the top without touching the edge.
export const SCROLL_OFFSET = 16

export function scrollToSection(id: string, offset: number = SCROLL_OFFSET) {
  const el = document.getElementById(id)
  if (!el) return
  const top = el.getBoundingClientRect().top + window.scrollY - offset
  window.scrollTo({ top: Math.max(top, 0), behavior: 'smooth' })
}
