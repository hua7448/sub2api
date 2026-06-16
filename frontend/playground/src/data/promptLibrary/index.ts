import type { PromptPresetSource, PromptPresetSummary } from './types'

export const PROMPT_LIBRARY_SOURCES: PromptPresetSource[] = [
  {
    "id": "youmind",
    "name": "YouMind OpenLab / Awesome GPT Image 2 Prompts",
    "url": "https://github.com/YouMind-OpenLab/awesome-gpt-image-2",
    "license": "CC BY 4.0",
    "licenseUrl": "https://creativecommons.org/licenses/by/4.0/",
    "commit": "7ae3665618f4181e2eb9e14b0d223f54a53b1478"
  },
  {
    "id": "davidwuw",
    "name": "davidwuw0811-boop / awesome-gpt-image2-prompts",
    "url": "https://github.com/davidwuw0811-boop/awesome-gpt-image2-prompts",
    "license": "MIT",
    "licenseUrl": "https://github.com/davidwuw0811-boop/awesome-gpt-image2-prompts",
    "commit": "228edb6341d978aa2572911adf7e7e147ebf95d3"
  },
  {
    "id": "freestylefly",
    "name": "freestylefly / awesome-gpt-image-2",
    "url": "https://github.com/freestylefly/awesome-gpt-image-2",
    "license": "MIT",
    "licenseUrl": "https://github.com/freestylefly/awesome-gpt-image-2/blob/main/LICENSE",
    "commit": "5ffa75c6e22e804b54a462841dc3d7035d8584ed"
  }
]

export const PROMPT_LIBRARY_CHUNK_IDS = [
  "chunk-0",
  "chunk-1",
  "chunk-2",
  "chunk-3",
  "chunk-4",
  "chunk-5"
] as const

export const PROMPT_LIBRARY_INDEX: PromptPresetSummary[] = [
  {
    "id": "freestylefly-1",
    "locale": "zh",
    "title": "信息图可视化设计",
    "category": "UI 与界面",
    "description": "Vertical 9:16 isometric cutaway infographic \"城市生命系统图谱 / Urban Metabolism Atlas\". Smart city from sky to bedrock: skyscrapers, streets, subwa",
    "promptPreview": "Vertical 9:16 isometric cutaway infographic \"城市生命系统图谱 / Urban Metabolism Atlas\". Smart city from sky to bedrock: skyscrapers, streets, subway, utility tunnels, ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "小红书号insight_express",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Infographic",
      "Tech",
      "Education",
      "Travel"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-330",
    "locale": "zh",
    "title": "月下美女直播画面",
    "category": "UI 与界面",
    "description": "生成一张直播间的图片，直播间氛围是月下美女跳舞的画面，直播间有很多人评论",
    "promptPreview": "生成一张直播间的图片，直播间氛围是月下美女跳舞的画面，直播间有很多人评论",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "苍何原创实测（公众号文章《我逆向了 329 条 GPT-Image2 提示词模板，全部开源！》）",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Creative"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-2",
    "locale": "zh",
    "title": "社媒界面截图",
    "category": "UI 与界面",
    "description": "画一张 X 的内容截图，深色模式，@OpenAI 蓝勾认证账号发推。 正文的中文内容： 今天想推荐一位很棒的 AI Builder：Ailln AI。 他持续在小红书分享 AI 工具、Agent 工作流、自动化实践和真实项目经验，把复杂的 AI 能力讲得清楚、实用、可落地。 如果",
    "promptPreview": "画一张 X 的内容截图，深色模式，@OpenAI 蓝勾认证账号发推。 正文的中文内容： 今天想推荐一位很棒的 AI Builder：Ailln AI。 他持续在小红书分享 AI 工具、Agent 工作流、自动化实践和真实项目经验，把复杂的 AI 能力讲得清楚、实用、可落地。 如果你正在关注 AI 产品、效率工具、个人自",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "小红书号4264014889",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Poster",
      "Tech",
      "Social"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-166",
    "locale": "zh",
    "title": "十二黄金圣斗士卡牌合集",
    "category": "其他案例",
    "description": "[中文] 生成圣斗士星矢12个黄金圣斗士的12宫格卡牌图片，每张卡牌上写上对应的中文名，每行4个，宽高比16:9。 [English] Generate a 12-grid card image of the 12 Gold Saints from Saint Seiya, wi",
    "promptPreview": "[中文] 生成圣斗士星矢12个黄金圣斗士的12宫格卡牌图片，每张卡牌上写上对应的中文名，每行4个，宽高比16:9。 [English] Generate a 12-grid card image of the 12 Gold Saints from Saint Seiya, with the corresponding",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@songguoxiansen",
    "tags": [
      "Other Use Cases",
      "Other Use Cases",
      "Tech"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-338",
    "locale": "zh",
    "title": "《赤壁怀古》长卷图",
    "category": "历史与古典主题",
    "description": "帮我生成一张《赤壁怀古》的长卷图，带整篇《赤壁赋》文字",
    "promptPreview": "帮我生成一张《赤壁怀古》的长卷图，带整篇《赤壁赋》文字",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "苍何原创实测（公众号文章《我逆向了 329 条 GPT-Image2 提示词模板，全部开源！》）",
    "tags": [
      "History & Classical Themes",
      "History",
      "Creative"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-354",
    "locale": "zh",
    "title": "Logo 与品牌身份系统提示词合集",
    "category": "品牌与标志",
    "description": "1. Logo概念生成提示词 你是一位拥有20年经验的顶级Logo设计师，为全球知名品牌设计过即时识别且深具意义的标志。 品牌名称：[你的品牌名] 行业：[你的行业] 品牌个性：[描述] 目标受众：[描述] 欣赏的视觉身份：[列举3个] 讨厌的视觉身份：[列举3个] 偏好风格：[",
    "promptPreview": "1. Logo概念生成提示词 你是一位拥有20年经验的顶级Logo设计师，为全球知名品牌设计过即时识别且深具意义的标志。 品牌名称：[你的品牌名] 行业：[你的行业] 品牌个性：[描述] 目标受众：[描述] 欣赏的视觉身份：[列举3个] 讨厌的视觉身份：[列举3个] 偏好风格：[如极简、大胆、几何、有机、复古、未来] ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@wanerfu",
    "tags": [
      "Brand & Logos",
      "Product",
      "Brand",
      "Tech",
      "Commerce"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-344",
    "locale": "en",
    "title": "NOIR 街头服饰 Campaign",
    "category": "品牌与标志",
    "description": "Create a premium, highly realistic 1:1 campaign poster for NOIR, a modern streetwear brand. Show one hero oversized hoodie as the main focus",
    "promptPreview": "Create a premium, highly realistic 1:1 campaign poster for NOIR, a modern streetwear brand. Show one hero oversized hoodie as the main focus against a gritty ur",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Daniel_adsss",
    "tags": [
      "Brand & Logos",
      "Poster",
      "Realistic",
      "Brand",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-362",
    "locale": "en",
    "title": "抹茶品牌触点系统视觉板",
    "category": "品牌与标志",
    "description": "Create a premium “Matcha Brand Touchpoint System” visual board for a modern lifestyle brand called: “MATCHA MODE” Build a full brand identit",
    "promptPreview": "Create a premium “Matcha Brand Touchpoint System” visual board for a modern lifestyle brand called: “MATCHA MODE” Build a full brand identity system, not a sing",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Preda2005",
    "tags": [
      "Brand & Logos",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-310",
    "locale": "zh",
    "title": "零食品牌技术分解图",
    "category": "品牌与标志",
    "description": "[中文] 创建一个 [SNACK] 的品牌技术信息图，结合产品的真实照片或照片级真实渲染，并将技术注释覆盖层直接置于其上。在纯白摄影棚背景上使用带有策略性 [BRAND COLOR] 点缀的黑色墨水风格线条画（建筑草图外观），包括： • 关键组件标签 • 显示结构、分层或内部设计",
    "promptPreview": "[中文] 创建一个 [SNACK] 的品牌技术信息图，结合产品的真实照片或照片级真实渲染，并将技术注释覆盖层直接置于其上。在纯白摄影棚背景上使用带有策略性 [BRAND COLOR] 点缀的黑色墨水风格线条画（建筑草图外观），包括： • 关键组件标签 • 显示结构、分层或内部设计的内部截面图 • 测量数据、尺寸和规格 ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@TechieBySA",
    "tags": [
      "Brand & Logos",
      "Infographic",
      "Realistic",
      "Illustration",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-370",
    "locale": "en",
    "title": "Crumple Chair 概念沙发研发板",
    "category": "商品与电商",
    "description": "Design Concept: The Crumple Chair Core Philosophy: Translating the \"controlled chaos\" of a tossed paper ball into a sculptural, high-comfort",
    "promptPreview": "Design Concept: The Crumple Chair Core Philosophy: Translating the \"controlled chaos\" of a tossed paper ball into a sculptural, high-comfort seating experience.",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@ShamsAmin56",
    "tags": [
      "Products & E-commerce",
      "UI",
      "Poster",
      "3D",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-365",
    "locale": "en",
    "title": "科学家收藏级玩具发布板",
    "category": "商品与电商",
    "description": "2x2 grid, do this for 4 famous scientists in history: Design a collector-grade launch visual for [TOY / FIGURE / DESIGNER OBJECT] shown in p",
    "promptPreview": "2x2 grid, do this for 4 famous scientists in history: Design a collector-grade launch visual for [TOY / FIGURE / DESIGNER OBJECT] shown in pristine hero form al",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Gdgtify",
    "tags": [
      "Products & E-commerce",
      "UI",
      "Product",
      "Brand",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-373",
    "locale": "zh",
    "title": "高端肉类海鲜品牌英雄图",
    "category": "商品与电商",
    "description": "一、品牌基础设定 品牌名称：[请填写，例如：PRIME STEAK / OCEAN PRIME] 品牌标语：[请填写，例如：Steakhouse Quality, Your Table / Restaurant Grade, Home Delivered] 主色调：[请填写，例如",
    "promptPreview": "一、品牌基础设定 品牌名称：[请填写，例如：PRIME STEAK / OCEAN PRIME] 品牌标语：[请填写，例如：Steakhouse Quality, Your Table / Restaurant Grade, Home Delivered] 主色调：[请填写，例如：黑金 / 深红+金 / 深蓝+银] 字",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@xpg0970",
    "tags": [
      "Products & E-commerce",
      "Brand",
      "Commerce",
      "Food"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-341",
    "locale": "en",
    "title": "AP Calculus 学习表信息图",
    "category": "图表与信息可视化",
    "description": "Please create a mathematical visualization infographic about \"[math concept / topic].\" The goal is to help the viewer intuitively understand",
    "promptPreview": "Please create a mathematical visualization infographic about \"[math concept / topic].\" The goal is to help the viewer intuitively understand what it is, why it ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@hqmank",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-375",
    "locale": "zh",
    "title": "古希腊三哲时间轴城市图",
    "category": "图表与信息可视化",
    "description": "二千五百年前，柏拉图，苏格拉底， 亚力士多德，坐在雅典街头聊天，聊出了世界文明史的源头。 背景可以加上他们聊天内容，按时间轴的走向，重叠在古希腊雅典的城市风光中。",
    "promptPreview": "二千五百年前，柏拉图，苏格拉底， 亚力士多德，坐在雅典街头聊天，聊出了世界文明史的源头。 背景可以加上他们聊天内容，按时间轴的走向，重叠在古希腊雅典的城市风光中。",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@ToroJushiAi",
    "tags": [
      "Charts & Infographics",
      "Charts",
      "Travel",
      "History"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-353",
    "locale": "zh",
    "title": "品牌口红推荐报告信息图",
    "category": "图表与信息可视化",
    "description": "一、系统角色 你是一个专业美妆顾问 + 人脸分析系统 + 品牌视觉设计系统。 你的任务是：基于用户上传自拍与指定口红品牌，生成一张具有品牌调性的“口红推荐报告信息结构图”。 二、输入参数 用户图像：{用户自拍} 品牌：{口红品牌，如 Dior / YSL / Armani / C",
    "promptPreview": "一、系统角色 你是一个专业美妆顾问 + 人脸分析系统 + 品牌视觉设计系统。 你的任务是：基于用户上传自拍与指定口红品牌，生成一张具有品牌调性的“口红推荐报告信息结构图”。 二、输入参数 用户图像：{用户自拍} 品牌：{口红品牌，如 Dior / YSL / Armani / Chanel / TF} 风格偏好（可选）",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Poster",
      "Commerce",
      "Food",
      "Story"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-361",
    "locale": "zh",
    "title": "手机爆炸拆解图",
    "category": "图表与信息可视化",
    "description": "Create a 3D Insane detailed exploded assembly drawing of [subject or object]",
    "promptPreview": "Create a 3D Insane detailed exploded assembly drawing of [subject or object]",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Ankit_patel211",
    "tags": [
      "Charts & Infographics",
      "3D",
      "Tech"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-360",
    "locale": "en",
    "title": "长发造型分析信息图",
    "category": "图表与信息可视化",
    "description": "Create a professional \"HAIRSTYLE ANALYSIS\" infographic with a different male model (the same face) having long, thick hair (6-10 inches), sl",
    "promptPreview": "Create a professional \"HAIRSTYLE ANALYSIS\" infographic with a different male model (the same face) having long, thick hair (6-10 inches), slightly wavy texture.",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Gemalpha_88",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-6",
    "locale": "zh",
    "title": "插画艺术创作图",
    "category": "插画与艺术",
    "description": "参考图是角色人设图，为参考图的少女绘制一副日系唯美奇幻风格插画。 【构图】这是一个宏大的中景日系奇幻插画构图，画面中心是完全保留了完整细节的可爱少女，她站立在无边的、如镜面般平滑的水面中心。天空呈现出高饱和度的粉紫与深蓝交织，一条耀眼的蓝色巨型流星划破天际，配合着边缘发光的瑰丽层",
    "promptPreview": "参考图是角色人设图，为参考图的少女绘制一副日系唯美奇幻风格插画。 【构图】这是一个宏大的中景日系奇幻插画构图，画面中心是完全保留了完整细节的可爱少女，她站立在无边的、如镜面般平滑的水面中心。天空呈现出高饱和度的粉紫与深蓝交织，一条耀眼的蓝色巨型流星划破天际，配合着边缘发光的瑰丽层云。女孩处于背光状态，形成一个暗调但依然",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "小红书号yi_xiao_jiu",
    "tags": [
      "Illustration & Art",
      "Illustration",
      "Character",
      "Creative"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-346",
    "locale": "zh",
    "title": "立体刺绣小鸟花枝",
    "category": "插画与艺术",
    "description": "精致立体刺绣风插画，浅浮雕纤维艺术效果，纯净「蚕丝白 + 奶白」底色，细腻丝线质感。画面为数只小鸟停在蜿蜒花枝上，周围点缀粉白、浅桃、珊瑚粉、淡金色花朵与叶片，构图轻盈雅致、留白充足。鸟儿羽毛以奶白、浅蓝、淡粉、浅金丝线刺绣表现，花枝纤细自然，花朵层层叠线，整体呈现高级手工刺绣、",
    "promptPreview": "精致立体刺绣风插画，浅浮雕纤维艺术效果，纯净「蚕丝白 + 奶白」底色，细腻丝线质感。画面为数只小鸟停在蜿蜒花枝上，周围点缀粉白、浅桃、珊瑚粉、淡金色花朵与叶片，构图轻盈雅致、留白充足。鸟儿羽毛以奶白、浅蓝、淡粉、浅金丝线刺绣表现，花枝纤细自然，花朵层层叠线，整体呈现高级手工刺绣、丝线堆绣、柔和光影、细节丰富、温柔清新的",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@dotey",
    "tags": [
      "Illustration & Art",
      "Illustration",
      "Creative"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-377",
    "locale": "en",
    "title": "樱花咖啡户外人像",
    "category": "摄影与写实",
    "description": "Edit the provided image while preserving the same face identity, shape, and facial features without altering age, ethnicity, or structure. M",
    "promptPreview": "Edit the provided image while preserving the same face identity, shape, and facial features without altering age, ethnicity, or structure. Maintain a calm, rela",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@xRahultripathi",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Brand",
      "Character",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-376",
    "locale": "en",
    "title": "泼洒抹茶街头手机照片",
    "category": "摄影与写实",
    "description": "A realistic vertical smartphone photo of a spilled green iced drink on outdoor stone pavement, a transparent disposable plastic cup lying on",
    "promptPreview": "A realistic vertical smartphone photo of a spilled green iced drink on outdoor stone pavement, a transparent disposable plastic cup lying on its side inside the",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Shinning1010",
    "tags": [
      "Photography & Realism",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-359",
    "locale": "en",
    "title": "水墨双重曝光人物海报",
    "category": "海报与排版",
    "description": "A cinematic character promotional poster of [SUBJECT], vertical composition (9:16), designed with a refined East-Asian ink aesthetic and hig",
    "promptPreview": "A cinematic character promotional poster of [SUBJECT], vertical composition (9:16), designed with a refined East-Asian ink aesthetic and high-end visual storyte",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Goodmanprotocol",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Character",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-350",
    "locale": "en",
    "title": "足球球员数据涂鸦海报",
    "category": "海报与排版",
    "description": "Create a scrapbook doodle-style football poster of [PLAYER_NAME]. AI auto-generates realistic career stats (club and national team). Main ph",
    "promptPreview": "Create a scrapbook doodle-style football poster of [PLAYER_NAME]. AI auto-generates realistic career stats (club and national team). Main photo: realistic, unch",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@ryanpp27",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Realistic",
      "Character",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "freestylefly-378",
    "locale": "en",
    "title": "高端 3D 收藏玩具头像",
    "category": "角色与人物",
    "description": "Transform the input photo into a high-end stylized 3D collectible figure. Large head, slightly exaggerated facial features while preserving ",
    "promptPreview": "Transform the input photo into a high-end stylized 3D collectible figure. Large head, slightly exaggerated facial features while preserving identity. Hyper-deta",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Genematicai",
    "tags": [
      "Characters & People",
      "Realistic",
      "Brand",
      "3D",
      "Tech",
      "Commerce"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-13491",
    "locale": "en",
    "title": "3D Stone Staircase Evolution Infographic",
    "category": "Featured Prompts",
    "description": "Transforms a flat evolutionary timeline into a realistic 3D stone staircase infographic with detailed organism renders and structured side p",
    "promptPreview": "{ \"type\": \"evolutionary timeline infographic\", \"instruction\": \"Using REFERENCE_0 as a structural base, transform the flat vector design into a highly realistic",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "知识猫图解",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-13467",
    "locale": "en",
    "title": "Anime Martial Arts Battle Illustration",
    "category": "Featured Prompts",
    "description": "Generates a dynamic, high-impact anime illustration of two female characters fighting in a traditional dojo with elemental energy effects.",
    "promptPreview": "An anime-style illustration of a {argument name=\"action type\" default=\"high-impact martial arts battle\"} between two young female fighters in a {argument name=\"",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "たねもみ 2.0 / Tanemomi Ver2.0",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25731",
    "locale": "en",
    "title": "Comic / Storyboard - Anime Post-Apocalyptic Scavenger Kids",
    "category": "Featured Prompts",
    "description": "A vertical manga-style character concept illustration of two scavenger children and a cat resting in a ruined sci-fi city.",
    "promptPreview": "Create a stylized anime / manga illustration of two young post-apocalyptic scavenger kids sitting side by side on a cracked concrete wall in a ruined city under",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Kōda",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25707",
    "locale": "en",
    "title": "Comic / Storyboard - Anime Travel Character Sheet Collage",
    "category": "Featured Prompts",
    "description": "Creates a playful manga-style reference-sheet collage from a character-and-car image, useful for turning a single illustration into multiple",
    "promptPreview": "Using the provided reference image as the character and car base, transform it into a lively hand-drawn Japanese anime character-reference collage on a white sk",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "あきを",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25614",
    "locale": "en",
    "title": "Comic / Storyboard - Bizarre Flat Cartoon Portrait",
    "category": "Featured Prompts",
    "description": "A stylized and surreal cartoon portrait prompt that features exaggerated geometric shapes, a playful design, and a flat graphic aesthetic.",
    "promptPreview": "Vertical bizarre flat cartoon portrait of {argument name=\"subject\" default=\"a girl from the attached photo\"} with a high geometric head shape, a long narrow nec",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Al-Shamus",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25492",
    "locale": "en",
    "title": "Comic / Storyboard - Dark Fantasy Anime Girl at Shrine",
    "category": "Featured Prompts",
    "description": "A detailed dark fantasy anime prompt featuring a girl at a mysterious shrine with lightning and spiritual talismans.",
    "promptPreview": "Masterpiece: Ultra-detailed {argument name=\"character type\" default=\"dark fantasy anime girl\"} based on the attached image, standing behind a partially open, an",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "HER19845",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25667",
    "locale": "en",
    "title": "Comic / Storyboard - Desert Cafe Escape Story",
    "category": "Featured Prompts",
    "description": "A narrative prompt for GPT Image 2 that generates a scene of a young boy resting at a quiet desert cafe, accompanied by a cat and a robot.",
    "promptPreview": "{argument name=\"title\" default=\"The Escaped Boy and a Cup in the Desert\"} A boy who escaped from a white ship is taking a short breath at a {argument name=\"loca",
    "sourceId": "youmind",
    "sourceLanguage": "JA",
    "sourceLabel": "もしもし@Aiart",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25669",
    "locale": "en",
    "title": "Comic / Storyboard - Devil Transformation Character Design",
    "category": "Featured Prompts",
    "description": "A sophisticated prompt for transforming an existing character into a dark fantasy devil version, preserving the original character's traits ",
    "promptPreview": "[Information Input Field] Character Name: [{argument name=\"character name\" default=\" \"}] (Required) Additional Notes: [{argument name=\"additional notes\" default",
    "sourceId": "youmind",
    "sourceLanguage": "JA",
    "sourceLabel": "muda22_Sora",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25700",
    "locale": "en",
    "title": "Comic / Storyboard - Moonlit Fantasy Alchemy Workshop",
    "category": "Featured Prompts",
    "description": "A detailed prompt for generating a dark fantasy alchemist laboratory scene with two robed figures, glowing potions, ancient books, and moonl",
    "promptPreview": "Create an ultra-detailed dark fantasy alchemist workshop scene at night, vertical 2:3 composition. Inside a cramped medieval stone-and-timber laboratory filled",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "おやぎ",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25699",
    "locale": "en",
    "title": "Comic / Storyboard - Overheated Office Lady at Noodle Shop",
    "category": "Featured Prompts",
    "description": "A detailed anime prompt for creating a sweaty office woman outside a Japanese shop advertising chilled hiyashi chuka noodles.",
    "promptPreview": "Create a vertical anime-style summer street illustration of a young office woman stopping outside a traditional Japanese noodle shop on a blazing hot day. The w",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "ヒー🎸✨",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25540",
    "locale": "en",
    "title": "Comic / Storyboard - Premium Pitch Deck Storyboard",
    "category": "Featured Prompts",
    "description": "A professional prompt for generating a 12-frame pitch deck storyboard with an editorial layout and a warm terracotta and sage color palette.",
    "promptPreview": "Create a high-end 4:3 pitch deck storyboard in 3x4 grid (12 frames), editorial layout, Headspace/Better Help premium campaign style, warm terracotta + soft sage",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "𝐌",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25708",
    "locale": "en",
    "title": "Comic / Storyboard - Retro Manga StackChan Concept Sheet",
    "category": "Featured Prompts",
    "description": "Creates a vintage Japanese manga-style character guide poster from StackChan robot references, ideal for concept art and customization inspi",
    "promptPreview": "Using REFERENCE_0 and REFERENCE_1 as the base design for the small StackChan wheeled robot, transform it into a cheerful retro Japanese manga/anime concept shee",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Yuto Flatmountain",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25620",
    "locale": "en",
    "title": "Comic / Storyboard - The Door to Another Life Surreal Fantasy Art",
    "category": "Featured Prompts",
    "description": "An ultra-cinematic prompt for a surreal fantasy scene of a man standing between a dark, rainy street and a radiant dream world island.",
    "promptPreview": "An ultra-cinematic surreal fantasy artwork called {argument name=\"artwork title\" default=\"The Door to Another Life\"}, a lonely handsome man standing on an empty",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Sheikh Sharik 2.0",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25670",
    "locale": "en",
    "title": "Comic / Storyboard - Vintage Fashion Caricature Sketch",
    "category": "Featured Prompts",
    "description": "A detailed prompt that transforms photos into adorable, fashion-inspired caricature sketches featuring ink lines and watercolors on vintage ",
    "promptPreview": "I transform the {argument name=\"subject\" default=\"people in my photographs\"} into adorable, {argument name=\"style\" default=\"fashion-inspired caricature sketches",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "AiRT🎥生成AI動画を創る人",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25630",
    "locale": "en",
    "title": "Comic / Storyboard - Watercolor Illustration of Guitar Player",
    "category": "Featured Prompts",
    "description": "A delicate watercolor and ink illustration of a young girl playing an acoustic guitar in a Japanese artbook style.",
    "promptPreview": "{ \"A delicate watercolor illustration of a young girl sitting cross-legged on the ground, softly playing a mint-green acoustic guitar. She has short dark hair w",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Maercih",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-14036",
    "locale": "en",
    "title": "E-commerce Live Stream UI Mockup",
    "category": "Featured Prompts",
    "description": "Generates a realistic social media live stream interface overlaying a portrait, featuring customizable chat messages, gift popups, and a pro",
    "promptPreview": "{ \"type\": \"live stream UI mockup\", \"subject\": { \"description\": \"portrait of {argument name=\\\"host name\\\" default=\\\"Elon Musk\\\"}, smiling, wearing a black t-shir",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "神经病不想好转",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25438",
    "locale": "en",
    "title": "E-commerce Main Image - Anonymous Tracksuit Lookbook Collage",
    "category": "Featured Prompts",
    "description": "Generates a realistic five-panel studio fashion collage for showcasing a women's hoodie-and-jogger tracksuit from multiple angles.",
    "promptPreview": "Goal: Create a studio fashion lookbook collage for a women's casual streetwear tracksuit, emphasizing fabric, fit, and multiple poses. Canvas: Wide horizontal 1",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Cherry 2.O",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-24875",
    "locale": "en",
    "title": "E-commerce Main Image - Editorial Fashion Photography JSON Prompt",
    "category": "Featured Prompts",
    "description": "A structured JSON prompt for fashion editorial photography featuring a monochromatic utilitarian chic jumpsuit.",
    "promptPreview": "subject\": {\"person\": \"woman\",\"appearance\": {\"hair\": \"dark brown, styled in a sleek, pulled-back bun\",\"pose\": \"standing, one leg lifted slightly with the knee be",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "NUSRAT",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25220",
    "locale": "en",
    "title": "E-commerce Main Image - Handmade Enamel Magnet Souvenir",
    "category": "Featured Prompts",
    "description": "Japanese translation of the prompt to create realistic enamel souvenir magnets with 3D depth and gold edges using Image 2.0.",
    "promptPreview": "Create a 3D handmade enamel souvenir magnet from a reference photo. Features recessed enamel surfaces, thin gold edges, soft studio light, an organic silhouette",
    "sourceId": "youmind",
    "sourceLanguage": "JA",
    "sourceLabel": "すん / 2つのアイ（AI・EYE）の専門家 / Codexであなたの仕事は自動化できる!",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25147",
    "locale": "en",
    "title": "E-commerce Main Image - High-Fashion Editorial Portrait Edit",
    "category": "Featured Prompts",
    "description": "A sophisticated prompt for editing reference photos into high-fashion studio portraits, focusing on color changes while preserving texture.",
    "promptPreview": "Edit the reference photo into a high-fashion studio portrait. Keep the same pose, facial structure, hairstyle, body proportions, lighting, and plain light gray",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Pixel by us ✨",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25663",
    "locale": "en",
    "title": "E-commerce Main Image - Korean Fashion Lookbook Grid Prompt",
    "category": "Featured Prompts",
    "description": "Generates a 2x3 grid editorial lookbook with strict identity consistency across six different outfits and poses.",
    "promptPreview": "Create a high-end Korean fashion lookbook featuring the same {argument name=\"subject\" default=\"Japanese woman\"} across all six looks, with strong identity consi",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "MatteoCruz",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-24870",
    "locale": "en",
    "title": "E-commerce Main Image - Korean Fashion Model Studio Portrait",
    "category": "Featured Prompts",
    "description": "A clean fashion editorial prompt for a Korean woman featuring specific hair styling and vibrant snakeskin accessories.",
    "promptPreview": "Create image featuring a young Korean woman with mid-length brown hair styled with {argument name=\"hair clip color\" default=\"lime green\"} hair clips on both sid",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "yusra.",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25618",
    "locale": "en",
    "title": "E-commerce Main Image - Korean Luxury Fashion Lookbook Grid",
    "category": "Featured Prompts",
    "description": "A comprehensive prompt for creating a 2x3 grid lookbook featuring a single model with consistent identity across six different high-end outf",
    "promptPreview": "The primary identity reference and display the exact same Japanese woman in all six panels with strict identity preservation. Keep identical facial structure, e",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "𝗦𝗮𝗻𝗶𝗮",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25252",
    "locale": "en",
    "title": "E-commerce Main Image - Lumina Mecha Figure Product Shot",
    "category": "Featured Prompts",
    "description": "A photorealistic premium toy photography prompt for generating a boxed anime-style armored android model kit named Lumina.",
    "promptPreview": "Create a photorealistic studio product shot of a futuristic collectible model kit figure named {argument name=\"character name\" default=\"LUMINA\"}. The scene show",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "カーブミラー",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-24902",
    "locale": "en",
    "title": "E-commerce Main Image - Luxury Candle Product Photography",
    "category": "Featured Prompts",
    "description": "Generates warm and atmospheric product images for candles and home fragrances, focusing on molten wax textures and cozy, high-end lifestyle ",
    "promptPreview": "Produce warm, atmospheric {argument name=\"product\" default=\"candle and home fragrance\"} product images that sell the feeling of luxury, comfort, and calm. This",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Kami AI",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25267",
    "locale": "en",
    "title": "E-commerce Main Image - Mauve Loungewear Bedroom Portrait",
    "category": "Featured Prompts",
    "description": "An ultra-realistic cozy bedroom lifestyle portrait of a young woman in mauve lounge clothing with a softly obscured face.",
    "promptPreview": "Create an ultra-realistic vertical lifestyle photograph of a {argument name=\"subject\" default=\"beautiful young woman with fair glowing skin\"} sitting relaxed on",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Alina Ai",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25047",
    "locale": "en",
    "title": "E-commerce Main Image - Modern Streetwear and Luxury T-Shirt Designs",
    "category": "Featured Prompts",
    "description": "A collection of professional prompts for designing premium, minimalist, and retro-inspired graphic t-shirts.",
    "promptPreview": "Create a high-end fashion product photo of a modern oversized streetwear T-shirt. The shirt features a large, professionally designed graphic print on the front",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Jack",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25036",
    "locale": "en",
    "title": "E-commerce Main Image - Premium Commercial Food Photography",
    "category": "Featured Prompts",
    "description": "A macro food photography prompt for a luxury advertising campaign, featuring high-end lighting and hyper-realistic textures.",
    "promptPreview": "Ultra-premium commercial food photography, hero shot of a {argument name=\"subject\" default=\"freshly baked crispy stuffed puff pastry (kachori)\"} placed on a smo",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Al-Shamus",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25493",
    "locale": "en",
    "title": "E-commerce Main Image - Strawberry Lava Mochi Explosion Photography",
    "category": "Featured Prompts",
    "description": "A high-end commercial food photography prompt featuring a strawberry mochi torn in half with molten filling and sparkling sugar crystals.",
    "promptPreview": "An ultra premium {argument name=\"dessert type\" default=\"Japanese strawberry mochi\"} suspended in mid air, violently torn into two halves at the exact moment of",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Snow",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-13515",
    "locale": "en",
    "title": "Illustrated City Food Map",
    "category": "Featured Prompts",
    "description": "Generates a hand-drawn, watercolor-style tourist map featuring numbered local food specialties, landmarks, and a legend.",
    "promptPreview": "{ \"type\": \"illustrated map infographic\", \"style\": \"{argument name=\\\"art style\\\" default=\\\"watercolor and ink hand-drawn illustration on vintage parchment\\\"}\", \"",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "皮皮特",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25643",
    "locale": "en",
    "title": "Infographic / Edu Visual - Camera Technical Exploded View",
    "category": "Featured Prompts",
    "description": "A precise technical illustration prompt for creating exploded-view diagrams of cameras with labeled internal components and a clean layout.",
    "promptPreview": "Detailed exploded-view diagram of a {argument name=\"camera model\" default=\"Sony A7 mirrorless camera\"}, with all internal components separated and clearly visib",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "PromptLab",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25647",
    "locale": "en",
    "title": "Infographic / Edu Visual - Cinematic Fashion Editorial Photography Composition Analysis",
    "category": "Featured Prompts",
    "description": "A detailed prompt for creating professional fashion editorial shots with an overlay of cinematic composition analysis, including grid lines ",
    "promptPreview": "Create cinematic fashion editorial photography in a {argument name=\"room style\" default=\"bright minimalist bedroom\"}, presented as a cinematographic composition",
    "sourceId": "youmind",
    "sourceLanguage": "ZH",
    "sourceLabel": "摆烂程序媛",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25421",
    "locale": "en",
    "title": "Infographic / Edu Visual - Codex Training Poster Workflow Infographic",
    "category": "Featured Prompts",
    "description": "A Japanese 16:9 infographic explaining a three-step Codex and Image 2.0 workflow for batch-generating four unified training seminar posters.",
    "promptPreview": "Goal: Create a clean Japanese presentation-style infographic showing how annual training seminar posters were batch-created with Codex and Image 2.0. Canvas: 16",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "てんそ｜医療× AI",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25414",
    "locale": "en",
    "title": "Infographic / Edu Visual - Dark TRPG Antagonist Character Sheet",
    "category": "Featured Prompts",
    "description": "Generates a classified dossier-style character design sheet for a modern occult TRPG antagonist with outfit details, turnarounds, portraits,",
    "promptPreview": "Goal: Create a dark tactical character reference sheet for a fictional TRPG antagonist, styled like a classified dossier and fashion design board. The character",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "安准",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25711",
    "locale": "en",
    "title": "Infographic / Edu Visual - DiT Parasite in Autoencoder Latent Diagram",
    "category": "Featured Prompts",
    "description": "A minimalist monochrome technical cartoon diagram showing a DiT creature living inside an autoencoder latent space for use in ML articles or",
    "promptPreview": "Create a minimalist black-and-white conceptual diagram on a wide 16:9 white canvas showing a diffusion-transformer creature parasitizing an autoencoder latent s",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "FA770",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25638",
    "locale": "en",
    "title": "Infographic / Edu Visual - Four Seasons Eye Macro Art",
    "category": "Featured Prompts",
    "description": "An artistic four-panel grid prompt representing the cycle of seasons through macro eye photography and thematic floral and lighting elements",
    "promptPreview": "Using an eye close-up as reference, create a 3:4 four-panel hyperrealistic eye macro, panels arranged top to bottom: Spring, Summer, Autumn, Winter. Panel 1 (Sp",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "PromptLab",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25673",
    "locale": "en",
    "title": "Infographic / Edu Visual - Gambling Aptitude Research Report",
    "category": "Featured Prompts",
    "description": "A comprehensive prompt for GPT Image 2 that generates a stylized 'gambling aptitude report' for a character, including stats and analysis ba",
    "promptPreview": "[Information Input Section] Character Name: [{argument name=\"character name\" default=\" \"}] *Required Supplementary Information: [{argument name=\"supplementary i",
    "sourceId": "youmind",
    "sourceLanguage": "JA",
    "sourceLabel": "muda22_Sora",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25378",
    "locale": "en",
    "title": "Infographic / Edu Visual - Global Beauty Standards Portraits",
    "category": "Featured Prompts",
    "description": "A comprehensive prompt for generating high-quality portraits that reflect the unique beauty standards and cultural aesthetics of twelve diff",
    "promptPreview": "Generate images of beautiful women and settings where they look best, as imagined by chatGPT and GPTimage2, for {argument name=\"target countries\" default=\"12 co",
    "sourceId": "youmind",
    "sourceLanguage": "JA",
    "sourceLabel": "陽仙堂",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25517",
    "locale": "en",
    "title": "Infographic / Edu Visual - Knitted Skyline Scarf Macro",
    "category": "Featured Prompts",
    "description": "A highly detailed macro photography prompt that visualizes a city skyline 3D-mapped into the threads of a knitted supporter scarf.",
    "promptPreview": "2x2 grid, do this for 4 dark horses for world cup 2026: 16:9, Anchor: [{argument name=\"city\" default=\"Country's famous city Skyline\"}] :: [Knitted Supporter Sca",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Gadgetify",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25389",
    "locale": "en",
    "title": "Infographic / Edu Visual - Master Potter Ceramic Infographic Grid",
    "category": "Featured Prompts",
    "description": "A complex prompt that creates a 2x2 grid of master potters merged with their ceramic materials and glaze chemistry in a cinematic studio set",
    "promptPreview": "2x2 grid, 16:9 do this for {argument name=\"number of potters\" default=\"4\"} famous master potters input: [{argument name=\"potter data\" default=\"four master potte",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Gadgetify",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25672",
    "locale": "en",
    "title": "Infographic / Edu Visual - MBTI Personality Character Card Poster",
    "category": "Featured Prompts",
    "description": "A comprehensive prompt for creating a high-end business magazine-style poster that visualizes MBTI personality types with detailed infograph",
    "promptPreview": "A vertical A4 poster in the style of a high-end Japanese business magazine cover. At the center is a character representing the '{argument name=\"personality typ",
    "sourceId": "youmind",
    "sourceLanguage": "JA",
    "sourceLabel": "荒井悠輔｜税理士×AI実務活用",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25398",
    "locale": "en",
    "title": "Infographic / Edu Visual - Pad Thai Cook Storyboard",
    "category": "Featured Prompts",
    "description": "Creates a professional and clean infographic storyboard poster with a 3D Pixar-style rendering, focusing on culinary themes and vibrant colo",
    "promptPreview": "Create a crisp, clean infographic storyboard poster for {argument name=\"dish\" default=\"THE PAD THAI COOK\"}. Wide 16:9 layout, white background, black borders, b",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "TechieSA",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25547",
    "locale": "en",
    "title": "Infographic / Edu Visual - Retro Robot Powering Illo with Codex",
    "category": "Featured Prompts",
    "description": "A whimsical ink-and-halftone illustration of a cute robot connecting an illo printer to a Codex subscription battery, suitable for developer",
    "promptPreview": "Create a playful retro black-and-white technical cartoon illustration on a warm off-white paper background, in a 16:9 wide canvas. Use thick rounded black ink o",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Trevin Chow",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25525",
    "locale": "en",
    "title": "Infographic / Edu Visual - Scientific Infographic for Pharmaceutical Research",
    "category": "Featured Prompts",
    "description": "A detailed scientific infographic prompt designed to visualize pharmaceutical concepts, mechanisms, and real-world research applications.",
    "promptPreview": "Create a detailed scientific infographic about {argument name=\"topic\" default=\"LIME Drug Design\"}, covering its core concepts, mechanisms, and real-world applic",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "PromptLab",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25420",
    "locale": "en",
    "title": "Infographic / Edu Visual - Three Japanese Note GPT Cover Proposals",
    "category": "Featured Prompts",
    "description": "Generates a three-panel comparison image of Japanese promotional cover design proposals for a note-specialized GPT service.",
    "promptPreview": "Goal: Create a wide comparison preview image showing exactly three proposed Japanese cover designs for a service called {argument name=\"service name\" default=\"n",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "みたかずお🌟ラッキーデザイナー",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-13983",
    "locale": "en",
    "title": "Momotaro Explainer Slide in Hybrid Style",
    "category": "Featured Prompts",
    "description": "A prompt that combines the simple, warm aesthetic of Irasutoya illustrations with the high-information density characteristic of Japanese go",
    "promptPreview": "Create an explanatory slide ({argument name=\"format\" default=\"ponchi-e diagram\"}) for {argument name=\"theme\" default=\"Momotaro\"} that fuses the gentle atmospher",
    "sourceId": "youmind",
    "sourceLanguage": "JA",
    "sourceLabel": "やまもん",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25631",
    "locale": "en",
    "title": "Product Marketing - 3D Tactical Crochet Art Diorama",
    "category": "Featured Prompts",
    "description": "A highly detailed prompt for creating a 3D yarn-based crochet diorama of a building with rich textile textures and pop-out effects.",
    "promptPreview": "A highly detailed 3D tactical crochet art diorama of a pastel blue, two story seaside shop from a straight on perspective, inspired by The building is crocheted",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Soulful Ai",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25725",
    "locale": "en",
    "title": "Product Marketing - Crystal Golem Card Comparison UI",
    "category": "Featured Prompts",
    "description": "Generates a side-by-side benchmark-style screenshot comparing two fantasy Crystal Golem trading card designs in a dark neon UI.",
    "promptPreview": "Goal: Create a dark comparison screenshot showing two fantasy trading-card designs side by side, as if from an image generation benchmark UI. Canvas: Landscape",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Nick",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25728",
    "locale": "en",
    "title": "Product Marketing - Editorial Riyad Mahrez Sports Poster",
    "category": "Featured Prompts",
    "description": "A bold Adidas-style football poster featuring a full-body athlete, oversized vertical typography, monochrome portrait collage, and red geome",
    "promptPreview": "Goal: Create a premium vertical sports poster design for {argument name=\"athlete name\" default=\"Riyad Mahrez\"}, styled like an Adidas football campaign, with a",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Ayyoub Ai",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25619",
    "locale": "en",
    "title": "Product Marketing - FIFA World Cup 2026 Sports Poster",
    "category": "Featured Prompts",
    "description": "A high-end cinematic sports poster prompt for the 2026 World Cup, featuring dynamic energy effects and customizable player and country detai",
    "promptPreview": "type: image_generation_prompt { \"title\": \"Universal FIFA World Cup 2026 Locked In Poster\", \"aspect_ratio\": \"9:16\", \"style\": \"ultra-realistic cinematic sports po",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Saul Goodman",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25683",
    "locale": "en",
    "title": "Product Marketing - FIFA World Cup Minimalist Poster",
    "category": "Featured Prompts",
    "description": "A high-end sports branding prompt for a minimalist World Cup poster with monochrome portraits, fragmented editorial effects, and premium typ",
    "promptPreview": "Create a premium FIFA World Cup poster design featuring a {argument name=\"player type\" default=\"legendary football superstar\"}. Minimalist white background with",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "K",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25709",
    "locale": "en",
    "title": "Product Marketing - Generative Bonsai Exhibition Poster",
    "category": "Featured Prompts",
    "description": "A minimalist Japanese-inspired exhibition poster featuring a photorealistic bonsai tree overlaid with generative geometry and refined editor",
    "promptPreview": "Goal: Create a minimalist Swiss-style exhibition poster for a generative nature bonsai show, centered on a realistic bonsai tree in a ceramic pot, combining bot",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Roku｜AI × BizDev",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25689",
    "locale": "en",
    "title": "Product Marketing - GPT Image 2 Reasoning Poster",
    "category": "Featured Prompts",
    "description": "A glossy 3D social poster showing a cube-brain icon and bold copy explaining GPT Image 2's plan-search-generate-verify workflow.",
    "promptPreview": "Goal: Create a polished vertical promotional poster for {argument name=\"product name\" default=\"GPT Image 2\"} explaining that the image model thinks before rende",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "The House of Curiosity",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25621",
    "locale": "en",
    "title": "Product Marketing - Luxury Skincare Sunscreen Ad Campaign",
    "category": "Featured Prompts",
    "description": "A professional-grade advertising prompt for a luxury sunscreen product, featuring cinematic tropical lighting and protective energy effects.",
    "promptPreview": "SUNSCREEN AD, “{argument name=\"headline\" default=\"THE INVISIBLE SHIELD\"}” Luxury skincare advertising masterpiece, a colossal premium sunscreen bottle standing",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Snow",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25727",
    "locale": "en",
    "title": "Product Marketing - Luxury Split Watch Campaign",
    "category": "Featured Prompts",
    "description": "A premium square fashion advertisement layout featuring two male models, architectural backdrops, luxury watch product shots, and elegant br",
    "promptPreview": "Create a luxurious square 1:1 high-end men's watch campaign advertisement for {argument name=\"brand name\" default=\"VANTORÉ\"}. The composition is split verticall",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "shah_zadii",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25729",
    "locale": "en",
    "title": "Product Marketing - Luxury VANTORÉ Watch Advertisement",
    "category": "Featured Prompts",
    "description": "A premium square fashion advertisement showing a suited male model and luxury Swiss wristwatch branding for use in high-end watch campaigns.",
    "promptPreview": "Goal: Create a luxury Swiss wristwatch advertisement for {argument name=\"brand name\" default=\"VANTORÉ\"}, featuring a single stylish male model in a premium fash",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "shah_zadii",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25616",
    "locale": "en",
    "title": "Product Marketing - Mechanical Wind-up Miniature World",
    "category": "Featured Prompts",
    "description": "An ultra-refined prompt for creating complex mechanical toy scenes with visible wind-up mechanisms and miniature narrative elements.",
    "promptPreview": "create a charming but ultra-refined scene centered on {argument name=\"subject\" default=\"wind-up toy / mechanical miniature world\"} where a tiny self-contained w",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Gadgetify",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25723",
    "locale": "en",
    "title": "Product Marketing - Oakley Skater Sunglasses Ad",
    "category": "Featured Prompts",
    "description": "A mixed-media eyewear brand ad combining a hand-drawn skater cartoon with realistic sunglasses and bold minimalist copy.",
    "promptPreview": "Create a minimalist vertical brand advertisement poster for {argument name=\"brand name\" default=\"OAKLEY\"}. On a clean white studio background, draw a black-and-",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Ilia Kitchenko",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25675",
    "locale": "en",
    "title": "Product Marketing - Product Design Catalog Layout",
    "category": "Featured Prompts",
    "description": "A highly detailed prompt for creating professional product design catalog pages, featuring a lifestyle hero shot and a technical specificati",
    "promptPreview": "Create a vertical 3:4 {argument name=\"page type\" default=\"product design catalog page\"} with a {argument name=\"background\" default=\"warm neutral paper-like back",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "PromptLab",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25686",
    "locale": "en",
    "title": "Profile / Avatar - Ancient Chinese Fantasy Character",
    "category": "Featured Prompts",
    "description": "A high-precision prompt for creating elegant 9:16 vertical portraits of women in an ancient Chinese fantasy aesthetic.",
    "promptPreview": "Generate a {argument name=\"composition\" default=\"9:16 vertical\"}, high-precision ancient Chinese fantasy-style refined portrait of a woman. The main subject is",
    "sourceId": "youmind",
    "sourceLanguage": "ZH",
    "sourceLabel": "李岳",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25644",
    "locale": "en",
    "title": "Profile / Avatar - Authentic Smartphone Night Selfie",
    "category": "Featured Prompts",
    "description": "Simulates a realistic, candid social media selfie captured indoors at night with smartphone camera aesthetics and low-light skin textures.",
    "promptPreview": "Ultra-realistic close-up selfie portrait of a {argument name=\"subject\" default=\"young woman\"} indoors at night, captured with a smartphone camera. Soft natural",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Jahan Zaib",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25633",
    "locale": "en",
    "title": "Profile / Avatar - Bizarre Geometric Cartoon Portrait",
    "category": "Featured Prompts",
    "description": "A surreal cartoon portrait prompt that transforms a reference photo into a geometric illustration with a stylized head shape and playful des",
    "promptPreview": "Vertical bizarre flat cartoon portrait of {argument name=\"subject\" default=\"a person\"} with a high geometric head shape, a long narrow neck, huge round eyes, a",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Saul Goodman",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25695",
    "locale": "en",
    "title": "Profile / Avatar - Censored Luxury Portrait",
    "category": "Featured Prompts",
    "description": "Generates a photorealistic close-up portrait of a stylish Black man in a sunlit luxury residence with the face obscured by a square censor b",
    "promptPreview": "Create an ultra-realistic close-up semi-portrait of {argument name=\"character description\" default=\"an exceptionally handsome Black man in his late 20s\"}, photo",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Finaltouch",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25698",
    "locale": "en",
    "title": "Profile / Avatar - Cinematic Xianxia Swordswoman Portrait",
    "category": "Featured Prompts",
    "description": "Generates a vertical Chinese fantasy drama-style close-up of a robed female cultivator holding a glowing emerald sword, useful for character",
    "promptPreview": "Create a vertical cinematic close-up portrait of {argument name=\"character name\" default=\"a graceful female immortal cultivator\"} in a Chinese xianxia fantasy d",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "𝟡𝟜 ᴾᴸᴬʸᶠᴼᴿᴳᴱ",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25681",
    "locale": "en",
    "title": "Profile / Avatar - Cyberpunk Neon Male Portrait",
    "category": "Featured Prompts",
    "description": "A high-contrast cinematic prompt for a cyberpunk-style male portrait featuring dramatic neon lighting, shallow depth of field, and magazine-",
    "promptPreview": "Cinematic close-up portrait of a man use image for face reference with same hair and light stubble, looking over his shoulder toward the camera with an intense,",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Duet | AI",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25611",
    "locale": "en",
    "title": "Profile / Avatar - Facial Identity Portrait for GPT Image 2",
    "category": "Featured Prompts",
    "description": "A specialized prompt designed to preserve exact facial identity while placing the subject in a lush outdoor garden setting with traditional ",
    "promptPreview": "Use the uploaded face reference image as the exact face reference with 100% facial identity preservation and full face lock. Do not change the face shape, eyes,",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Mahnoor Fatima",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25612",
    "locale": "en",
    "title": "Profile / Avatar - Fan Photo with Cristiano Ronaldo",
    "category": "Featured Prompts",
    "description": "Generates a highly realistic sports photography style image of a user standing alongside Cristiano Ronaldo at a World Cup stadium while pres",
    "promptPreview": "Use the uploaded image as the ONLY facial identity reference with STRICT identity preservation enabled at maximum strength. Do not modify, beautify, reshape, or",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Harboris",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25685",
    "locale": "en",
    "title": "Profile / Avatar - Japanese Anime Girl Illustration",
    "category": "Featured Prompts",
    "description": "A simple prompt to generate beautiful characters in a classic Japanese anime art style using the high-performance rendering of GPT Image 2.",
    "promptPreview": "Draw a {argument name=\"subject\" default=\"beautiful girl\"} in {argument name=\"style\" default=\"Japanese anime style\"}",
    "sourceId": "youmind",
    "sourceLanguage": "JA",
    "sourceLabel": "濃縮ビンプル@個人ゲーム制作者",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25705",
    "locale": "en",
    "title": "Profile / Avatar - Moody Portrait With Privacy Block",
    "category": "Featured Prompts",
    "description": "A cinematic square close-up portrait with the subject’s face hidden by a large brown gradient privacy block, useful for anonymous avatar or ",
    "promptPreview": "Create a square, cinematic close-up portrait of {argument name=\"subject\" default=\"an androgynous young adult with tousled dark hair\"} in a dim indoor setting, c",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "かッッき",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25715",
    "locale": "en",
    "title": "Profile / Avatar - Moonlit Anime Woman with Peony",
    "category": "Featured Prompts",
    "description": "A detailed fantasy anime illustration prompt for creating a melancholic moonlit woman holding a pink peony beneath a starry night sky.",
    "promptPreview": "Create a vertical 2:3 fantasy anime-style illustration of one sorrowful young woman in left-facing side profile, standing under a star-filled midnight sky. She",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "星_AIart",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25617",
    "locale": "en",
    "title": "Profile / Avatar - Photo to Realistic Sticker Pack",
    "category": "Featured Prompts",
    "description": "A prompt that converts a single photo into a set of realistic digital stickers with various expressions and poses, perfect for messaging app",
    "promptPreview": "Turn this image into a {argument name=\"style\" default=\"realistic sticker pack\"} featuring multiple expressions, poses, emotions, and reactions of the subject. A",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Virena",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25702",
    "locale": "en",
    "title": "Profile / Avatar - Sunlit Maid Sisters Hug",
    "category": "Featured Prompts",
    "description": "A warm anime-style portrait of two maid sisters embracing in a softly lit room, useful for character illustration prompts.",
    "promptPreview": "Create a warm, high-detail anime illustration of two maid sisters hugging in a sunlit room. The scene shows exactly 2 characters: an older sister on the right,",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Matsubon",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25628",
    "locale": "en",
    "title": "Profile / Avatar - Vibrant Watercolor Anime Portrait",
    "category": "Featured Prompts",
    "description": "A warm and colorful watercolor anime-style portrait of a handsome man in an orange shirt, set against a botanical background.",
    "promptPreview": "Cute watercolor anime-style portrait illustration of a handsome young man with soft wavy dark hair and oversized transparent round glasses, resting his chin on",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Aegon",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25622",
    "locale": "en",
    "title": "Profile / Avatar - Watercolor Anime Male Portrait",
    "category": "Featured Prompts",
    "description": "A dreamy and whimsical watercolor portrait of a handsome young man with glasses, featuring soft pastel tones and delicate botanical elements",
    "promptPreview": "Cute watercolor anime-style portrait illustration of a handsome young man(same as in uploaded image) with soft wavy {argument name=\"hair color\" default=\"black\"}",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Orion",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25639",
    "locale": "en",
    "title": "Profile / Avatar - Watercolor Anime Male Portrait",
    "category": "Featured Prompts",
    "description": "A whimsical prompt for creating soft watercolor anime-style character illustrations with botanical backgrounds and hand-painted textures.",
    "promptPreview": "Cute watercolor anime-style portrait illustration of a {argument name=\"character\" default=\"handsome young man\"}(same as in uploaded image) with soft wavy black",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Jahan Zaib",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25615",
    "locale": "en",
    "title": "Profile / Avatar - Woman with Kitten Close-up",
    "category": "Featured Prompts",
    "description": "A charming eye-level close-up of a young woman playfully hiding behind a fluffy kitten, set against a lush greenery background.",
    "promptPreview": "A close-up, eye-level shot of a young woman with long, wavy brown hair and front bangs holding up a {argument name=\"kitten color\" default=\"fluffy, light grey an",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "yusra.",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25652",
    "locale": "en",
    "title": "Social Media Post - Ancient Abandoned Library Fantasy Portrait",
    "category": "Featured Prompts",
    "description": "A cinematic fantasy portrait of a person in an abandoned library surrounded by floating books and golden sunlight.",
    "promptPreview": "Use my uploaded face image as the primary identity reference. Preserve my exact facial identity, face shape, hairstyle, hair texture, beard style, skin tone, ey",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Professor",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25653",
    "locale": "en",
    "title": "Social Media Post - B&W Fashion Casting Contact Sheet",
    "category": "Featured Prompts",
    "description": "A monochrome fashion grid prompt that creates a 2x2 grid of different expressions and angles for an editorial test shoot look.",
    "promptPreview": "Black-and-white fashion casting contact sheet of a woman with {argument name=\"hair style\" default=\"short wet curls\"} glistening under studio lights, arranged in",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Virena",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25654",
    "locale": "en",
    "title": "Social Media Post - Celestial Observatory Explorer",
    "category": "Featured Prompts",
    "description": "A highly detailed fantasy prompt of an adventurer in a floating observatory surrounded by brass instruments and cats.",
    "promptPreview": "Use the uploaded image as the exact facial identity reference. Preserve facial structure, jawline, hairstyle, skin texture, facial proportions, expression reali",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "WeWant Mars",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25720",
    "locale": "en",
    "title": "Social Media Post - Cinematic Pet Photo Enhancer",
    "category": "Featured Prompts",
    "description": "Transforms a low-resolution pet photo into a sharp, photorealistic cinematic-quality image while preserving identity and composition.",
    "promptPreview": "Using the provided reference image, enhance and upscale it into an ultra-premium cinematic-quality photo. Preserve the subject's identity, pose, framing, and ov",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Geloria",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25660",
    "locale": "en",
    "title": "Social Media Post - Cinematic Urban Fashion Editorial Portrait",
    "category": "Featured Prompts",
    "description": "Transforms an uploaded portrait into a stunning cinematic lifestyle editorial featuring luxury automotive elements and golden-hour lighting.",
    "promptPreview": "Using the uploaded image and without changing face create a stunning cinematic lifestyle portrait of the {argument name=\"subject\" default=\"woman\"} seated confid",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "simeon-sanai",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25693",
    "locale": "en",
    "title": "Social Media Post - French Femme Fatale Poetry Poster",
    "category": "Featured Prompts",
    "description": "A wide vintage-style art poster pairs a censored melodrama painting of two women with handwritten poem text, credits, links, and a QR-code g",
    "promptPreview": "Goal: Create a horizontal social-media art poster combining a vintage French melodrama painting with a black poetry-and-credits panel. Canvas: Wide 16:9 composi",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Michael H. Lester",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25642",
    "locale": "en",
    "title": "Social Media Post - Golden Hour Close-up Portrait",
    "category": "Featured Prompts",
    "description": "An intimate cinematic close-up prompt focused on warm amber lighting, natural skin details, and a moody film aesthetic.",
    "promptPreview": "Cinematic close-up portrait of a {argument name=\"subject\" default=\"beautiful woman with short wavy brunette bob hair\"} partially covering one eye, soft natural",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Jahan Zaib",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25724",
    "locale": "en",
    "title": "Social Media Post - GPT Image 2 Space Porthole",
    "category": "Featured Prompts",
    "description": "A photorealistic sci-fi promotional scene showing a spacecraft porthole view of a planet embossed with GPT Image 2 text.",
    "promptPreview": "Create a cinematic photorealistic sci-fi view from inside a spacecraft cockpit or round observation window, looking out through a thick dark metal circular port",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "sadeghTM",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25696",
    "locale": "en",
    "title": "Social Media Post - Magical Sunrise Alpine Valley",
    "category": "Featured Prompts",
    "description": "A photorealistic cinematic wilderness scene with deer drinking in a misty alpine river valley at sunrise, ideal for nature, travel, or fanta",
    "promptPreview": "Create a cinematic, ultra-realistic nature landscape at magical sunrise: a pristine alpine valley with a shallow winding river in the foreground, framed by dens",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Nahe Kola",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25688",
    "locale": "en",
    "title": "Social Media Post - Miniature Girl on Giant Smartphone Desk",
    "category": "Featured Prompts",
    "description": "A cozy pastel miniature workspace scene featuring a tiny stylish girl sitting on a giant smartphone, useful for whimsical lifestyle and soci",
    "promptPreview": "Create a warm, whimsical photorealistic miniature scene in a cozy pastel study workspace: a tiny cute young South Asian woman, {argument name=\"character name\" d",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Mehwish kiran",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25730",
    "locale": "en",
    "title": "Social Media Post - Noir Luxury Fashion Editorial",
    "category": "Featured Prompts",
    "description": "A cinematic fashion photograph prompt for creating an anonymous model in a tailored suit inside a glamorous nighttime skyscraper interior.",
    "promptPreview": "Create a cinematic high-fashion editorial photograph of a poised adult woman in an elegant {argument name=\"suit color\" default=\"olive green\"} tailored pantsuit,",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Zephyra Leigh",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25694",
    "locale": "en",
    "title": "Social Media Post - Paparazzi Arrival at Luxury Hotel",
    "category": "Featured Prompts",
    "description": "A realistic celebrity street-photo scene showing an elegant woman exiting a luxury hotel amid security, photographers, and a parked black se",
    "promptPreview": "Create a realistic paparazzi-style street photograph of {argument name=\"character name\" default=\"an elegant woman\"} stepping out of the Four Seasons hotel entra",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Nadia Ai",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25636",
    "locale": "en",
    "title": "Social Media Post - Realistic Photo with Cartoon Avatar",
    "category": "Featured Prompts",
    "description": "A creative prompt that places a realistic person next to their own high-quality 3D cartoon character counterpart based on a reference photo.",
    "promptPreview": "Create a high-quality realistic photo featuring the {argument name=\"subject\" default=\"uploaded person\"} standing in a stylish pose. Beside them, place a {argume",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Smiling Khan",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25613",
    "locale": "en",
    "title": "Social Media Post - Skateboarding Girl in Traditional Sari",
    "category": "Featured Prompts",
    "description": "A cinematic documentary-style photography prompt featuring a young girl in a traditional sari performing at a skatepark, creating a bold con",
    "promptPreview": "A cinematic documentary style photograph of a young girl {argument name=\"action\" default=\"skateboarding\"} at an outdoor skatepark while wearing a {argument name",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "BMX",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25678",
    "locale": "en",
    "title": "Social Media Post - Three Frame Portrait Collage Prompt",
    "category": "Featured Prompts",
    "description": "A high-quality prompt for generating a realistic three-frame collage of a woman featuring warm golden hour lighting and soft aesthetic detai",
    "promptPreview": "Ultra-realistic aesthetic 3-frame portrait collage of a {argument name=\"subject\" default=\"beautiful young woman\"}, soft natural makeup, glowing fair skin, rosy",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "🐣",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25687",
    "locale": "en",
    "title": "Social Media Post - Twilight Park Stair Portrait",
    "category": "Featured Prompts",
    "description": "A cinematic photorealistic portrait prompt for creating a young woman seated on old stone stairs beside a glowing park lamp at dusk.",
    "promptPreview": "Create a cinematic hyper-realistic portrait of a {argument name=\"character age and gender\" default=\"21-year-old woman\"} sitting casually on weathered stone stai",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Alina Ai",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25658",
    "locale": "en",
    "title": "Social Media Post - Victorian Woman Riding Vintage Bicycle",
    "category": "Featured Prompts",
    "description": "A cinematic and ultra-realistic portrait of a woman in a Victorian dress riding a vintage bicycle through a sun-drenched autumn park.",
    "promptPreview": "An ultra-realistic, cinematic portrait of a {argument name=\"subject\" default=\"beautiful young woman\"} riding a {argument name=\"bicycle style\" default=\"vintage b",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Avelyrah",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-13460",
    "locale": "en",
    "title": "VR Headset Exploded View Poster",
    "category": "Featured Prompts",
    "description": "Generates a high-tech exploded view diagram of a VR headset with detailed component callouts and promotional text.",
    "promptPreview": "{ \"type\": \"exploded view product diagram poster\", \"subject\": \"VR headset\", \"style\": \"clean high-tech 3D render, studio lighting, glowing accents\", \"background\":",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "wory＠ホッピング中",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25415",
    "locale": "en",
    "title": "YouTube Thumbnail - Anime Idol Advertising Takeover City",
    "category": "Featured Prompts",
    "description": "A wide cinematic prompt for generating a rainy neon city filled with coordinated billboard ads of one anime idol character.",
    "promptPreview": "Create a cinematic wide 16:9 anime-style nighttime cityscape showing a full “advertising takeover” of a futuristic Tokyo-like entertainment district by one idol",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "オオトリヒカリ🐦‍🔥",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-23633",
    "locale": "en",
    "title": "YouTube Thumbnail - Anime Protagonist and Companion Generator",
    "category": "Featured Prompts",
    "description": "An extensive image-to-image prompt for creating a new protagonist character that stylistically matches and narratively complements a referen",
    "promptPreview": "Use the provided image as the base reference for the character. If the character in the reference image is an important figure from a non-existent masterpiece a",
    "sourceId": "youmind",
    "sourceLanguage": "JA",
    "sourceLabel": "Eris Create Lab",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-24216",
    "locale": "en",
    "title": "YouTube Thumbnail - Anime Typhoon Live News Reporter",
    "category": "Featured Prompts",
    "description": "Generates a Japanese live TV news scene featuring a pink idol-style anime reporter covering a typhoon on a stormy Kobe waterfront.",
    "promptPreview": "Goal: Create a dramatic live television news still of an anime typhoon report in Kobe, Japan, combining a realistic stormy waterfront background with a cute ani",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Kiki",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-24676",
    "locale": "en",
    "title": "YouTube Thumbnail - ChatGPT Foundations Product Visual",
    "category": "Featured Prompts",
    "description": "A premium square course-thumbnail graphic showing a glowing AI-generated product mockup surrounded by prompt-engineering terminal windows.",
    "promptPreview": "Create a polished square promotional hero graphic for an online class about using ChatGPT for commercial product visuals. Canvas: 1:1 aspect ratio, dark navy-to",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "DDS Vibe Academy",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25209",
    "locale": "en",
    "title": "YouTube Thumbnail - Cinematic Film Genre Color Grading Comparison",
    "category": "Featured Prompts",
    "description": "A detailed prompt for creating a cinematic vertical poster that compares five different film color grading styles for the same urban street ",
    "promptPreview": "Create a {argument name=\"aspect\" default=\"vertical\"} cinematic tutorial comparison poster showing the same {argument name=\"subject\" default=\"urban street scene\"",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "BMX",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-24193",
    "locale": "en",
    "title": "YouTube Thumbnail - Correct Face in Existing Thumbnail",
    "category": "Featured Prompts",
    "description": "Generates a revised YouTube thumbnail by preserving the full design while correcting the subject’s face using a second reference image.",
    "promptPreview": "Using REFERENCE_1 as the current YouTube-style thumbnail and REFERENCE_0 as the identity reference, regenerate the thumbnail with the same layout and compositio",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "無職51歳 初めて実験物語",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-23456",
    "locale": "en",
    "title": "YouTube Thumbnail - Futuristic Cyberpunk Creative Leader Poster",
    "category": "Featured Prompts",
    "description": "A prompt for creating a high-end cinematic cyberpunk poster featuring a creative leader in a luxury office with holographic elements.",
    "promptPreview": "Use the provided reference image of me for facial resemblance only. Do not copy the original pose, clothing, hairstyle, background, or lighting. Transform me in",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Zeeshan",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-24521",
    "locale": "en",
    "title": "YouTube Thumbnail - GPT Image 2 Cooking Prompt Banner",
    "category": "Featured Prompts",
    "description": "A cute Japanese promotional banner showing four cooking-image prompt examples for recipe cards, food photos, photo books, and movie-poster v",
    "promptPreview": "Goal: Create a wide Japanese promotional thumbnail/banner for {argument name=\"headline topic\" default=\"GPT Image 2 cooking image prompts\"}, designed like a cute",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Maki@Sunwood AI Labs.",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-24005",
    "locale": "en",
    "title": "YouTube Thumbnail - GPT Image 2.0 Japanese Thumbnail",
    "category": "Featured Prompts",
    "description": "Generates a bold Japanese YouTube thumbnail showing a before-and-after AI photo enhancement scene for GPT Image 2.0 content.",
    "promptPreview": "Create a high-impact Japanese YouTube thumbnail about AI image generation improvement. Canvas: 16:9 horizontal thumbnail, photorealistic with bold creator-thumb",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "たつみん∣AIで仕組化するPdM",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-23954",
    "locale": "en",
    "title": "YouTube Thumbnail - Gritty Neo-Noir Tactical Portrait",
    "category": "Featured Prompts",
    "description": "A cinematic prompt for a high-contrast tactical portrait in an abandoned garage with neo-noir lighting.",
    "promptPreview": "**Shot Type:** Wide Cinematic Medium-Full Shot **Subject & Expression:** {argument name=\"subject\" default=\"A realistic portrait featuring the exact man from the",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Dilshad Hussain",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-24004",
    "locale": "en",
    "title": "YouTube Thumbnail - Japanese Lawson Milk Can Thumbnail",
    "category": "Featured Prompts",
    "description": "A vibrant anime-style Japanese YouTube thumbnail featuring a presenter, bold headline text, a Lawson Day badge, and a milk-can visual gag.",
    "promptPreview": "Create a bright Japanese YouTube thumbnail in a high-energy variety-show style about the urban legend that the blue Lawson sign is actually a milk can. Use a 16",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "癒色えも(イシキ•エモ)",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-25204",
    "locale": "en",
    "title": "YouTube Thumbnail - Minecraft Death Screen Portrait",
    "category": "Featured Prompts",
    "description": "A realistic winter photography prompt that blends a photographic portrait with Minecraft game elements like the 'You Died' screen and pixela",
    "promptPreview": "Use 100% of the face and body from my reference image. Create a cinematic winter scene with a side face close-up view, inspired by a Minecraft death screen. The",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "CHAse",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-23277",
    "locale": "en",
    "title": "YouTube Thumbnail - Miniature Set Style Portrait",
    "category": "Featured Prompts",
    "description": "A unique cinematic high-angle shot of a character on a scooter within a meticulously crafted miniature city set.",
    "promptPreview": "A cinematic wide-angle shot features uploaded face as reference, riding a {argument name=\"vehicle\" default=\"vintage light-green scooter\"}, wearing a {argument n",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Dilshad Hussain",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-24208",
    "locale": "en",
    "title": "YouTube Thumbnail - NeoNippon Imperial Throne Hall",
    "category": "Featured Prompts",
    "description": "A detailed prompt for generating a cinematic fictional NeoNippon imperial palace interior with Japanese-inspired gilded architecture and war",
    "promptPreview": "Create a grand cinematic interior scene in a fictional aesthetic called {argument name=\"aesthetic name\" default=\"NeoNippon\"}: a vast imperial throne hall blendi",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "SakuraFan99",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-24503",
    "locale": "en",
    "title": "YouTube Thumbnail - Shadow Sync Anime Key Visual",
    "category": "Featured Prompts",
    "description": "A cinematic split-world anime poster featuring two blurred-face protagonists between a bright utopia and a neon cyberpunk city.",
    "promptPreview": "Create a cinematic anime key visual for a sci-fi fantasy project titled {argument name=\"title text\" default=\"SHADOW SYNC\"}. Wide 16:9 poster composition, ultra-",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Shadow Sync",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-24972",
    "locale": "en",
    "title": "YouTube Thumbnail - Shogi Battlefield Video Thumbnail",
    "category": "Featured Prompts",
    "description": "A cinematic thumbnail prompt for a fantasy shogi battlefield scene starring two anthropomorphic dogs with award-style social video overlays.",
    "promptPreview": "Goal: Create a cinematic social media video thumbnail / award post image showing a fantasy shogi match transforming into a battlefield, with a dark PixVerse-sty",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "pink shih tzu ponta",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-23782",
    "locale": "en",
    "title": "YouTube Thumbnail - Skyscraper Climbing Action Shot",
    "category": "Featured Prompts",
    "description": "An ultra-realistic, high-angle cinematic shot of a subject climbing a glass skyscraper with precise reference mapping.",
    "promptPreview": "ultra-realistic IMAX-grade cinematic action shot, extreme high-angle perspective looking down from a massive height, 9:16 vertical aspect ratio, a {argument nam",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Junaid Mohana",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-24824",
    "locale": "en",
    "title": "YouTube Thumbnail - Summon Dash Roblox Fantasy Thumbnail",
    "category": "Featured Prompts",
    "description": "Generates a dynamic Roblox-style fantasy game thumbnail with a running knight, glowing portal, collectible crystals, and bold promotional te",
    "promptPreview": "Goal: Create a high-energy Roblox fantasy game thumbnail for {argument name=\"game title\" default=\"SUMMON DASH\"}, styled like a premium mobile RPG advertisement",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "HONEBITO（映像クリエイター / 生成AI）",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-23596",
    "locale": "en",
    "title": "YouTube Thumbnail - Vibrant High-Contrast Digital Collage",
    "category": "Featured Prompts",
    "description": "A dynamic digital collage featuring multiple poses of a young woman on a bright red gradient background with modern typography.",
    "promptPreview": "A high-contrast, energetic digital collage on a vibrant solid {argument name=\"background color\" default=\"red gradient\"} background. Large uppercase text across",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Avelyrah",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-en-24964",
    "locale": "en",
    "title": "YouTube Thumbnail - Zombie Shoppers Point Sale Cover",
    "category": "Featured Prompts",
    "description": "A photorealistic Japanese blog cover scene of zombie-like shoppers crowding a store entrance while holding sale-point smartphones.",
    "promptPreview": "Create a cinematic photorealistic blog cover image showing a chaotic crowd of zombie-like shoppers pressed against the glass entrance of a Japanese retail store",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "おやぎ",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-13491",
    "locale": "zh",
    "title": "3D 石阶演化信息图",
    "category": "精选提示词",
    "description": "将平面的演化时间轴转化为逼真的 3D 石阶信息图，包含精细的生物渲染图和结构化的侧边栏。",
    "promptPreview": "{ \"type\": \"演化时间轴信息图\", \"instruction\": \"以 REFERENCE_0 作为结构基础，将平面矢量设计转化为高度逼真的 3D 信息图。将平滑的坡道替换为独特的石阶，并将所有生物升级为照片级真实的 3D 模型。\", \"style\": { \"background\": \"{argument na",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "知识猫图解",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-13460",
    "locale": "zh",
    "title": "VR 头显爆炸视图海报",
    "category": "精选提示词",
    "description": "生成一张高科技 VR 头显爆炸视图，包含详细的组件标注和宣传文案。",
    "promptPreview": "{ \"type\": \"产品爆炸视图海报\", \"subject\": \"VR 头显\", \"style\": \"简洁的高科技 3D 渲染，摄影棚灯光，发光装饰\", \"background\": \"{argument name=\\\"background color\\\" default=\\\"柔和的紫蓝色渐变\\\"}\", \"header",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "wory＠ホッピング中",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-24964",
    "locale": "zh",
    "title": "YouTube 缩略图 - “丧尸”购物狂潮促销封面",
    "category": "精选提示词",
    "description": "一张照片级写实的日本博客封面场景，描绘了如丧尸般的购物者挤在商店门口，手中举着显示促销点的智能手机。",
    "promptPreview": "创作一张电影级照片写实的博客封面图，展示在阳光明媚的日子里，一群如丧尸般混乱的购物者挤在日本零售店玻璃入口处的场景。摄像机位于店内，透过自动滑动玻璃门向外拍摄，采用 16:9 宽屏构图，低位略带广角视角，强烈的逆光，窗外是蓝天和朵朵白云，玻璃上有倒影，金属门框，侧窗贴有海报，右下角堆放着一叠红色购物篮。约 14 名脸色",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "おやぎ",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-24676",
    "locale": "zh",
    "title": "YouTube 缩略图 - ChatGPT 基础课程产品视觉图",
    "category": "精选提示词",
    "description": "一张优质的正方形课程缩略图，展示了由 AI 生成的发光产品模型，周围环绕着提示词工程终端窗口。",
    "promptPreview": "为一门关于使用 ChatGPT 进行商业产品视觉设计的在线课程，制作一张精美的正方形宣传主图。画布：1:1 比例，深海军蓝至炭灰色的渐变背景，带有微妙的摄影棚地面反射和柔和的暗角效果。构图中心悬浮着一件高级空白 {argument name=\"product type\" default=\"T 恤\"}，渲染为逼真的 3D",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "DDS Vibe Academy",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-24521",
    "locale": "zh",
    "title": "YouTube 缩略图 - GPT Image 2 烹饪提示词横幅",
    "category": "精选提示词",
    "description": "一个可爱的日式宣传横幅，展示了四个用于食谱卡片、美食摄影、写真集和电影海报视觉效果的烹饪图像提示词示例。",
    "promptPreview": "目标：为 {argument name=\"headline topic\" default=\"GPT Image 2 烹饪图像提示词\"} 创建一个宽幅日式宣传缩略图/横幅，设计风格类似于可爱的手工食谱拼贴画。 画布：横向 16:9 比例的社交媒体横幅，采用暖色调奶油纸背景，带有细腻的织物纹理、圆角贴纸形状、柔和的阴影，配",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Maki@Sunwood AI Labs.",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-24005",
    "locale": "zh",
    "title": "YouTube 缩略图 - GPT Image 2.0 日语缩略图",
    "category": "精选提示词",
    "description": "为 GPT Image 2.0 内容生成一张醒目的日语 YouTube 缩略图，展示 AI 照片增强前后的对比场景。",
    "promptPreview": "制作一张极具视觉冲击力的日语 YouTube 缩略图，主题为 AI 图像生成效果提升。画布：16:9 横向缩略图，照片级真实感，采用醒目的创作者缩略图排版。场景：温馨且柔和模糊的咖啡馆内景，背景有吊灯和光斑，前景为木质桌面。右侧展示 1 位年轻日本女性，身穿白色蕾丝衬衫，头发盘成优雅的松散发髻，侧身朝左坐着；她的脸部被",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "たつみん∣AIで仕組化するPdM",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-25204",
    "locale": "zh",
    "title": "YouTube 缩略图 - Minecraft 死亡界面肖像",
    "category": "精选提示词",
    "description": "一个写实的冬季摄影提示词，将人像摄影与 Minecraft 游戏元素（如“You Died”界面和像素化物品）巧妙融合。",
    "promptPreview": "100% 使用我参考图中的面部和身体。创作一个电影感的冬季场景，采用侧脸特写视角，灵感源自 Minecraft 死亡界面。人物仰面躺在雪地中，身穿 {argument name=\"clothing\" default=\"休闲冬装\"}，姿态自然写实，配以冷色调的环境光，色彩略微去饱和，雪地上留有清晰的脚印。在人物周围环绕或",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "CHAse",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-24208",
    "locale": "zh",
    "title": "YouTube 缩略图 - NeoNippon 帝国宝座大厅",
    "category": "精选提示词",
    "description": "一份详细的提示词，用于生成具有电影质感的虚构 NeoNippon 帝国宫殿室内场景，融合了日式风格的镀金建筑与温暖的体积光效果。",
    "promptPreview": "创建一个宏大的电影级室内场景，采用名为 {argument name=\"aesthetic name\" default=\"NeoNippon\"} 的虚构美学风格：一个广阔的帝国宝座大厅，将传统的日本神社和宫殿建筑与未来主义的奢华感融为一体。构图采用从地面视角拍摄的 4:3 广角镜头，对称且庄严，中心是一个高耸的礼仪台，",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "SakuraFan99",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-24503",
    "locale": "zh",
    "title": "YouTube 缩略图 - Shadow Sync 动画主视觉图",
    "category": "精选提示词",
    "description": "一张电影感十足的动画海报，展现了分裂的世界，两位面部模糊的主角站在明亮的乌托邦与霓虹赛博朋克城市之间。",
    "promptPreview": "创作一张电影感十足的动画主视觉图，用于名为 {argument name=\"title text\" default=\"SHADOW SYNC\"} 的科幻奇幻项目。采用 16:9 宽屏海报构图，超高细节，分裂世界对比：左半部分是明亮的乌托邦未来城市，白昼下有白色水晶摩天大楼、蓝天、太阳耀斑、绿植、水道、流畅的建筑线条和金",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Shadow Sync",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-24824",
    "locale": "zh",
    "title": "YouTube 缩略图 - Summon Dash Roblox 奇幻风格缩略图",
    "category": "精选提示词",
    "description": "生成一张动态的 Roblox 风格奇幻游戏缩略图，包含奔跑的骑士、发光的传送门、可收集的水晶以及醒目的宣传文字。",
    "promptPreview": "目标：为 {argument name=\"game title\" default=\"SUMMON DASH\"} 创建一张高能量的 Roblox 奇幻游戏缩略图，采用类似高级移动端 RPG 广告的风格，并配以醒目的漫画动作图形。 画布：16:9 宽屏缩略图，1200x675 分辨率，色彩鲜艳且饱和度高，旨在确保在 You",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "HONEBITO（映像クリエイター / 生成AI）",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-24193",
    "locale": "zh",
    "title": "YouTube 缩略图 - 修正现有缩略图中的人脸",
    "category": "精选提示词",
    "description": "通过使用第二张参考图修正主体面部，同时保留完整设计，从而生成修订后的 YouTube 缩略图。",
    "promptPreview": "以 REFERENCE_1 作为当前的 YouTube 风格缩略图，以 REFERENCE_0 作为身份参考，在保持相同布局和构图的前提下重新生成缩略图，但仅替换男性的面部，使其与 REFERENCE_0 中的人物一致。保持现有的姿势、指向手指、低音吉他、房间背景、灯光、高对比度缩略图风格以及所有日语排版不变。精确保留",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "無職51歳 初めて実験物語",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-23596",
    "locale": "zh",
    "title": "YouTube 缩略图 - 充满活力的强对比度数字拼贴画",
    "category": "精选提示词",
    "description": "一幅充满动感的数字拼贴画，以明亮的红色渐变背景为特色，展示了一位年轻女性的多种姿态，并配有现代排版设计。",
    "promptPreview": "一幅高对比度、充满活力的数字拼贴画，背景为鲜艳的 {argument name=\"background color\" default=\"红色渐变\"}。顶部的大号大写文字显示为“{argument name=\"branding\" default=\"AvelyrahnAI\"}”，采用简洁现代的字体。构图展示了同一位年轻女性",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Avelyrah",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-23633",
    "locale": "zh",
    "title": "YouTube 缩略图 - 动漫主角与伙伴生成器",
    "category": "精选提示词",
    "description": "一个详尽的图生图提示词，用于创作一个在 90 年代复古动漫风格上与参考角色在视觉和叙事上相辅相成的新主角。",
    "promptPreview": "使用提供的图像作为角色参考基准。如果参考图像中的角色是某部虚构的动漫、漫画或游戏杰作中的重要人物，请生成一个“最合适的主角”，要求该主角能与参考角色并肩作战、共同旅行、深度参与其命运，并使他们的故事成为可能。主角应是自然衍生自参考角色的存在。主角并非现有角色，必须是全新的设计。然而，他们必须具备压倒性的魅力和存在感，让",
    "sourceId": "youmind",
    "sourceLanguage": "JA",
    "sourceLabel": "Eris Create Lab",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-25415",
    "locale": "zh",
    "title": "YouTube 缩略图 - 动漫偶像广告霸屏城市",
    "category": "精选提示词",
    "description": "一个宽画幅电影感提示词，用于生成充满霓虹雨夜的城市，且所有广告牌均展示同一位动漫偶像角色。",
    "promptPreview": "创作一张 16:9 电影感宽画幅动漫风格夜景图，展现一个类似未来东京的娱乐区被同一位偶像角色“广告霸屏”的场景，{argument name=\"character name\" default=\"Hikari\"}。视角从前景的高处行人阳台俯瞰，下方是雨中霓虹闪烁的广场，挤满了微小的行人，街道湿润反光，玻璃建筑、列车、空中",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "オオトリヒカリ🐦‍🔥",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-24216",
    "locale": "zh",
    "title": "YouTube 缩略图 - 动漫台风现场新闻记者",
    "category": "精选提示词",
    "description": "生成一个日本电视直播新闻场景，画面中一位粉色偶像风格的动漫记者正在神户风暴肆虐的海滨报道台风。",
    "promptPreview": "目标：创作一张极具戏剧性的电视直播新闻剧照，展示动漫风格的台风报道。画面结合了写实的海滨风暴背景与前景中合成的可爱动漫记者。 画布：16:9 横向广播画幅，黄昏或傍晚时分写实的雨中海滨城市场景，灰蓝色风暴光影，湿润的反射路面，汹涌的海浪，密集的雨线，强风，乌云，远处的城市天际线和港口塔楼。 主体：一位可爱的金发动漫少女",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Kiki",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-24972",
    "locale": "zh",
    "title": "YouTube 缩略图 - 将棋战场视频缩略图",
    "category": "精选提示词",
    "description": "一个电影感十足的奇幻将棋战场场景缩略图提示词，主角为两只拟人化狗狗，并带有奖项风格的社交视频叠加层。",
    "promptPreview": "目标：创建一个电影感社交媒体视频缩略图 / 奖项发布图片，展示一场转化为战场的奇幻将棋对局，采用深色 PixVerse 风格的视频卡片呈现。 画布：600×483 像素的横向卡片，黑色背景，圆角设计。上方 82% 为电影级静态图像；下方为黑色长条区域，带有大号标题。使用温暖的火光、深邃的阴影、浅景深、戏剧性的电影布光以",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "pink shih tzu ponta",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-23277",
    "locale": "zh",
    "title": "YouTube 缩略图 - 微缩模型风格肖像",
    "category": "精选提示词",
    "description": "一张独特的电影感高角度镜头，展现了精心制作的微缩城市模型中，一名骑着滑板车的角色。",
    "promptPreview": "一张电影感广角镜头，以上传的人脸作为参考，骑着 {argument name=\"vehicle\" default=\"复古浅绿色滑板车\"}，穿着 {argument name=\"attire\" default=\"深色西装和红色领带\"}，表情严肃地直视镜头；滑板车位于一个风格化的地图状表面上，周围有道路和蜿蜒的蓝色路径，并",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Dilshad Hussain",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-23782",
    "locale": "zh",
    "title": "YouTube 缩略图 - 摩天大楼攀爬动作镜头",
    "category": "精选提示词",
    "description": "一张超写实、高角度的电影级镜头，展示主体在玻璃摩天大楼上攀爬，并带有精确的参考映射。",
    "promptPreview": "超写实 IMAX 级电影动作镜头，从极高处俯瞰的极限高角度视角，9:16 竖屏比例，一位 {argument name=\"subject\" default=\"男性或女性\"} 匹配精确的参考面部和发型，正在攀爬现代玻璃摩天大楼的侧面，身体在反光的幕墙上动态伸展，位于画面右侧，一条腿抬起抓紧玻璃，另一条腿弯曲支撑，双手紧贴",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Junaid Mohana",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-24004",
    "locale": "zh",
    "title": "YouTube 缩略图 - 日本 Lawson 牛奶罐缩略图",
    "category": "精选提示词",
    "description": "一张充满活力的动漫风格日本 YouTube 缩略图，包含一位主持人、醒目的标题文字、Lawson 日徽章以及牛奶罐视觉梗。",
    "promptPreview": "制作一张明亮的日本 YouTube 缩略图，采用高能综艺节目风格，主题是关于“蓝色 Lawson 招牌其实是牛奶罐”的都市传说。使用 16:9 横向画布，采用光泽感动漫/Vtuber 缩略图设计，背景为温暖模糊的便利店或城市灯光，添加厚重的白色轮廓、闪光、放射状速度线、半调网点、漫画飞溅效果和强烈的投影。在右侧放置一位",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "癒色えも(イシキ•エモ)",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-23456",
    "locale": "zh",
    "title": "YouTube 缩略图 - 未来赛博朋克创意领袖海报",
    "category": "精选提示词",
    "description": "用于创作一张高端电影感赛博朋克海报的提示词，展现一位身处奢华办公室、带有全息元素的创意领袖。",
    "promptPreview": "仅参考所提供的本人照片以获取面部相似度。请勿复制原图的姿势、服装、发型、背景或光影。将我塑造成一位强大的未来主义创意领袖 / 赛博朋克远见者。创作一张高端电影感赛博朋克海报，主角是一位自信的 {argument name=\"ethnicity\" default=\"南亚裔\"} 男性，坐在黑色真皮行政椅上，身着全黑奢华装束",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Zeeshan",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-25209",
    "locale": "zh",
    "title": "YouTube 缩略图 - 电影级影片类型调色对比",
    "category": "精选提示词",
    "description": "一份详细的提示词，用于创作一张电影级竖版海报，对比同一城市街景的五种不同电影调色风格。",
    "promptPreview": "创作一张 {argument name=\"aspect\" default=\"vertical\"} 的电影级教程对比海报，展示同一个 {argument name=\"subject\" default=\"城市街景\"} 在五个水平面板中重复出现，每个面板展示一种不同的电影类型调色风格。场景为日落时分繁忙的小城市街道，包含汽车",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "BMX",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-23954",
    "locale": "zh",
    "title": "YouTube 缩略图 - 硬核新黑色电影风格战术肖像",
    "category": "精选提示词",
    "description": "一个用于生成高对比度战术肖像的电影级提示词，场景设定在废弃车库中，采用新黑色电影风格布光。",
    "promptPreview": "**拍摄类型：** 宽画幅电影级中全景镜头 **主体与表情：** {argument name=\"subject\" default=\"一张写实肖像，主体为上传参考图中的男士\"}。他保持着冷静、近乎掠食者般的自信。 **场景与氛围：** {argument name=\"setting\" default=\"夜晚，雨水湿滑的",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Dilshad Hussain",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-0"
  },
  {
    "id": "youmind-zh-25611",
    "locale": "zh",
    "title": "个人资料 / 头像 - GPT Image 2 面部特征保持肖像",
    "category": "精选提示词",
    "description": "这是一个专门设计的提示词，旨在保留精确的面部特征，同时将主体置于郁郁葱葱的户外花园背景中，并身着传统花卉服饰。",
    "promptPreview": "请使用上传的面部参考图作为精确的面部参考，实现 100% 的面部特征保持和全脸锁定。请勿更改脸型、眼睛、眉毛、鼻子、嘴唇、肤色、发型、发际线、表情、年龄或任何独特的面部特征。仅替换面部为上传的参考图，同时保持其他所有元素与原始场景完全一致。 一位美丽的年轻女性坐在户外郁郁葱葱的花园庭院中，身着 {argument na",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Mahnoor Fatima",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25612",
    "locale": "zh",
    "title": "个人资料 / 头像 - 与 Cristiano Ronaldo 的合影",
    "category": "精选提示词",
    "description": "生成一张高度逼真的体育摄影风格照片，让用户与 Cristiano Ronaldo 在世界杯球场并肩合影，同时严格保留用户的面部特征。",
    "promptPreview": "请将上传的图片作为唯一的面部特征参考，并开启最高强度的身份保留功能。请勿修改、美化、重塑或重构我的面部结构。请精准保留我在参考图中的眼睛、鼻子、下颌线、笑容形状、肤色、发际线、发型、面部比例以及整体身份特征。 请创作一张超逼真的高质量照片，画面内容为我与 {argument name=\"celebrity\" defau",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Harboris",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25686",
    "locale": "zh",
    "title": "个人资料 / 头像 - 中国古风奇幻角色",
    "category": "精选提示词",
    "description": "一个高精度提示词，用于创作优雅的 9:16 竖屏中国古风奇幻女性肖像。",
    "promptPreview": "生成一张 {argument name=\"composition\" default=\"9:16 竖屏\"}、高精度的中国古风奇幻风格精致女性肖像。主体为一名 {argument name=\"age\" default=\"20–28\"} 岁左右的年轻东亚女性。",
    "sourceId": "youmind",
    "sourceLanguage": "ZH",
    "sourceLabel": "李岳",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25615",
    "locale": "zh",
    "title": "个人资料 / 头像 - 女子与小猫特写",
    "category": "精选提示词",
    "description": "一张迷人的平视特写照片，一位年轻女子俏皮地躲在毛茸茸的小猫身后，背景是郁郁葱葱的绿植。",
    "promptPreview": "这是一张平视特写镜头，一位留着棕色波浪长发和齐刘海的年轻女子将一只 {argument name=\"kitten color\" default=\"毛茸茸的浅灰白色小猫\"} 举在脸前，小猫身后仅露出她的眼睛、鼻子和前额。小猫有着又大又好奇的绿色眼睛、明显的胡须以及额头上深色的虎斑纹，正警觉地直视镜头。女子身穿鲜艳的红色上",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "yusra.",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25705",
    "locale": "zh",
    "title": "个人资料 / 头像 - 带有隐私遮挡的忧郁肖像",
    "category": "精选提示词",
    "description": "一张电影质感的方形特写肖像，主体面部被一个大型棕色渐变隐私遮挡块遮住，适用于匿名头像或概念图像。",
    "promptPreview": "创作一张方形、电影质感的特写肖像，主体为 {argument name=\"subject\" default=\"一位留着凌乱深色头发的中性化年轻人\"}，场景设定在昏暗的室内，进行紧凑裁剪，使头部和肩部填满画面。人物穿着一件 {argument name=\"clothing\" default=\"深炭灰色针织毛衣\"}，一只手",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "かッッき",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25633",
    "locale": "zh",
    "title": "个人资料 / 头像 - 怪诞几何卡通肖像",
    "category": "精选提示词",
    "description": "一个超现实的卡通肖像提示词，可将参考照片转换为具有风格化头部形状和趣味设计的几何插画。",
    "promptPreview": "垂直怪诞扁平卡通肖像，主角为 {argument name=\"subject\" default=\"一个人\"}，拥有高耸的几何头部形状、细长的脖子、巨大的圆眼睛、小巧的嘴巴和淡定的笑容，身穿 {argument name=\"clothing\" default=\"色彩鲜艳的街头服饰\"}，头上戴着 {argument nam",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Saul Goodman",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25685",
    "locale": "zh",
    "title": "个人资料 / 头像 - 日系动漫少女插画",
    "category": "精选提示词",
    "description": "一个简单的提示词，利用 GPT Image 2 的高性能渲染能力，生成经典日系动漫风格的精美角色。",
    "promptPreview": "绘制一个 {argument name=\"subject\" default=\"美丽的女孩\"}，风格为 {argument name=\"style\" default=\"日系动漫风格\"}",
    "sourceId": "youmind",
    "sourceLanguage": "JA",
    "sourceLabel": "濃縮ビンプル@個人ゲーム制作者",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25715",
    "locale": "zh",
    "title": "个人资料 / 头像 - 月光下的动漫少女与牡丹",
    "category": "精选提示词",
    "description": "一个精细的奇幻动漫插画提示词，用于创作一位在星空下手持粉色牡丹、神情忧郁的月光少女。",
    "promptPreview": "创作一张 2:3 竖构图的奇幻动漫风格插画，描绘一位神情忧郁的年轻女性，侧脸向左，伫立在繁星点点的午夜天空下。她拥有一头极长的飘逸 {argument name=\"hair color\" default=\"深蓝黑色\"} 秀发，发丝流淌在画面中，细节丰富；她的面部被一个位于中心、边缘柔和的灰褐色矩形模糊遮盖。她身着一件华",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "星_AIart",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25622",
    "locale": "zh",
    "title": "个人资料 / 头像 - 水彩动漫男性肖像",
    "category": "精选提示词",
    "description": "一幅梦幻且充满奇思妙想的水彩肖像，描绘了一位戴眼镜的英俊青年，采用柔和的粉彩色调和精致的植物元素。",
    "promptPreview": "可爱的水彩动漫风格肖像插画，描绘了一位英俊的年轻男子（与上传图片一致），留着柔和的波浪状 {argument name=\"hair color\" default=\"黑色\"} 头发，戴着超大号透明圆框眼镜，单手托腮，面带温柔微笑。拥有明亮闪烁且富有表现力的双眼，淡淡的红晕，柔和的粉彩肤色，精致的墨线勾勒，充满奇思妙想的水",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Orion",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25639",
    "locale": "zh",
    "title": "个人资料 / 头像 - 水彩动漫男性肖像",
    "category": "精选提示词",
    "description": "一个用于创作柔和水彩动漫风格角色插画的奇幻提示词，包含植物背景和手绘纹理。",
    "promptPreview": "可爱的水彩动漫风格肖像插画，主角为 {argument name=\"character\" default=\"英俊的年轻男子\"}（与上传图片一致），留着柔软的黑色波浪卷发，戴着 {argument name=\"glasses type\" default=\"超大透明圆框眼镜\"}，单手托腮，面带温柔微笑。拥有明亮闪烁的动人双",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Jahan Zaib",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25617",
    "locale": "zh",
    "title": "个人资料 / 头像 - 照片转写实贴纸包",
    "category": "精选提示词",
    "description": "一个可以将单张照片转换为一系列写实数字贴纸的提示词，包含多种表情和姿势，非常适合在即时通讯软件中使用。",
    "promptPreview": "将此图像转换为 {argument name=\"style\" default=\"写实贴纸包\"}，包含主体多种表情、姿势、情绪和反应。将贴纸排列在干净的版面上，以便在即时通讯软件和数字交流中使用。",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Virena",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25628",
    "locale": "zh",
    "title": "个人资料 / 头像 - 生动的水彩动漫肖像",
    "category": "精选提示词",
    "description": "一幅温暖且色彩丰富的水彩动漫风格肖像，描绘了一位身穿橙色衬衫的英俊男子，背景为植物元素。",
    "promptPreview": "可爱的水彩动漫风格肖像插画，描绘了一位留着柔软波浪卷深色头发、戴着超大透明圆框眼镜的英俊青年，他单手托腮，带着温暖柔和的微笑。拥有明亮闪烁且富有表现力的双眼、淡淡的红晕、柔和的粉彩肤色、精致的墨水草图轮廓，以及奇幻的水彩飞溅和晕染效果。身穿一件色彩鲜艳的 {argument name=\"shirt color\" def",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Aegon",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25713",
    "locale": "zh",
    "title": "个人资料 / 头像 - 电影院大厅里的哥特风动漫少女",
    "category": "精选提示词",
    "description": "一张写实的竖版时尚快照，描绘了一位哥特风动漫少女在日本电影院大厅的电影排片表下摆拍的场景。",
    "promptPreview": "创作一张写实的竖版手机摄影照片，场景位于夜晚的现代电影院大厅，展示了一位名为 {argument name=\"character name\" default=\"Emilia\"} 的动漫风格年轻女性，她以四分之三侧面站立在巨大的数字“NOW SHOWING”项目下方。她留着非常长的 {argument name=\"hai",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "レティシア・ノエル",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25644",
    "locale": "zh",
    "title": "个人资料 / 头像 - 真实感智能手机夜间自拍",
    "category": "精选提示词",
    "description": "模拟在室内夜间环境下，通过智能手机摄像头拍摄的真实、随性的社交媒体自拍效果，呈现低光环境下的皮肤质感。",
    "promptPreview": "一张超写实的近景自拍肖像，拍摄对象为 {argument name=\"subject\" default=\"年轻女性\"}，场景设定在室内夜间，由智能手机摄像头拍摄。皮肤质感柔和自然，肤色细腻，脸颊泛着淡淡的红晕，双唇粉嫩有光泽，表情俏皮，舌尖轻触上唇。明亮深邃的双眸略微侧向镜头。中长深棕色秀发配有轻盈的法式刘海，修饰着脸",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Jahan Zaib",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25681",
    "locale": "zh",
    "title": "个人资料 / 头像 - 赛博朋克霓虹男性肖像",
    "category": "精选提示词",
    "description": "一个高对比度的电影感提示词，用于创作赛博朋克风格的男性肖像，包含戏剧性的霓虹灯光、浅景深效果以及杂志级的摄影质感。",
    "promptPreview": "电影感男性特写肖像，使用人物图像作为面部参考，保留相同的发型和轻微胡茬，转头看向镜头，表情自信且强烈，佩戴同款透明椭圆形墨镜。{argument name=\"lighting style\" default=\"戏剧性霓虹灯光\"}，以 {argument name=\"key color\" default=\"深蓝色\"} 作为",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Duet | AI",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25702",
    "locale": "zh",
    "title": "个人资料 / 头像 - 阳光下的女仆姐妹拥抱",
    "category": "精选提示词",
    "description": "一幅温馨的动漫风格肖像，描绘了两位女仆姐妹在柔光房间中拥抱的场景，适用于角色插画提示词。",
    "promptPreview": "创作一幅温馨且高细节的动漫插画，描绘两位女仆姐妹在阳光明媚的房间里拥抱。画面中恰好有 2 个角色：右侧是姐姐，身材较高，身穿经典的黑色长袖女仆装，配有白色荷叶边围裙、白色褶皱女仆发带、领口处有黑色蝴蝶结，留着 {argument name=\"older sister hair color\" default=\"暖棕色\"}",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Matsubon",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25619",
    "locale": "zh",
    "title": "产品营销 - 2026 年 FIFA 世界杯体育海报",
    "category": "精选提示词",
    "description": "一款专为 2026 年世界杯打造的高端电影感体育海报提示词，包含动态能量特效，并支持自定义球员和国家信息。",
    "promptPreview": "type: image_generation_prompt { \"title\": \"2026 年 FIFA 世界杯“全力以赴”通用海报\", \"aspect_ratio\": \"9:16\", \"style\": \"超写实电影感体育海报\", \"quality\": \"16K UHD 极致写实\", \"subject\": { \"id",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Saul Goodman",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25631",
    "locale": "zh",
    "title": "产品营销 - 3D 战术钩针艺术立体模型",
    "category": "精选提示词",
    "description": "一个高细节的提示词，用于创作一个 3D 纱线钩针立体模型，呈现具有丰富纺织纹理和立体突出效果的建筑。",
    "promptPreview": "一个高度精细的 3D 战术钩针艺术立体模型，呈现了一座正面视角的淡蓝色两层海边商店。建筑由整齐、紧密的浅蓝色纱线针脚钩织而成，配有桃粉色的纱线门框和窗框。左侧是一朵巨大、蓬松、极具质感的积雨云，采用厚实的白色圈圈纱制成，营造出 3D 突出效果，背景是浅蓝色的钩针天空。下方，海洋由深浅蓝色纱线的波浪状横排组成，沙滩则由质",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Soulful Ai",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25725",
    "locale": "zh",
    "title": "产品营销 - Crystal Golem 卡牌对比 UI",
    "category": "精选提示词",
    "description": "生成一张并排的基准测试风格截图，对比两种深色霓虹 UI 风格的奇幻 Crystal Golem 集换式卡牌设计。",
    "promptPreview": "目标：创建一个深色对比截图，展示两张并排的奇幻集换式卡牌设计，呈现出图像生成基准测试 UI 的效果。 画布：16:9 横向构图，黑色/深炭灰色背景，带有细细的霓虹绿色轮廓。两个高大的垂直面板填满画面，左右各一，中间有窄缝隔开。每个面板包含一个可收藏的怪物卡牌预览，下方是一个小型状态/控制区域。 布局：使用 2 个主面板",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Nick",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25683",
    "locale": "zh",
    "title": "产品营销 - FIFA 世界杯极简海报",
    "category": "精选提示词",
    "description": "一款高端体育品牌推广提示词，用于创作极简风格的世界杯海报，包含单色肖像、碎片化编辑效果及高级排版。",
    "promptPreview": "创作一张高级的 FIFA 世界杯海报设计，主角为 {argument name=\"player type\" default=\"传奇足球巨星\"}。极简白色背景，画面主体为巨大的单色特写肖像。肖像被分割成若干水平切片，切片间留有细长的白色缝隙，营造出摩登的碎片化编辑效果。左侧采用大型垂直粗体排版，以 {argument n",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "K",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25689",
    "locale": "zh",
    "title": "产品营销 - GPT Image 2 推理海报",
    "category": "精选提示词",
    "description": "一张光泽感 3D 社交媒体海报，展示了一个立方体大脑图标，并配有醒目的文案，解释了 GPT Image 2 的“规划-搜索-生成-验证”工作流程。",
    "promptPreview": "目标：为 {argument name=\"product name\" default=\"GPT Image 2\"} 创建一张精致的竖版宣传海报，通过俏皮的 3D 玩具风格视觉效果，解释该图像模型在渲染前会进行思考。 画布：4:5 竖版社交媒体海报，深钴蓝色背景，带有微妙的径向光晕和柔和的晕影，呈现高分辨率光泽感 3D",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "The House of Curiosity",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25723",
    "locale": "zh",
    "title": "产品营销 - Oakley 滑板风格太阳镜广告",
    "category": "精选提示词",
    "description": "一个混合媒介眼镜品牌广告，将手绘滑板卡通形象与写实的太阳镜及大胆的极简主义文案相结合。",
    "promptPreview": "为 {argument name=\"brand name\" default=\"OAKLEY\"} 创建一张极简主义风格的竖版品牌广告海报。在干净的白色摄影棚背景上，绘制一个黑白水墨风格的卡通滑板少年，他有着刺猬头、大头、坏笑的表情、交叉的双臂，穿着 T 恤、长短裤、袜子和运动鞋。他自信地站立，一只脚着地，另一只脚踩在倾斜",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Ilia Kitchenko",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25728",
    "locale": "zh",
    "title": "产品营销 - Riyad Mahrez 体育海报编辑设计",
    "category": "精选提示词",
    "description": "一款大胆的 Adidas 风格足球海报，包含全身运动员形象、超大垂直排版、单色肖像拼贴以及红色几何装饰元素。",
    "promptPreview": "目标：为 {argument name=\"athlete name\" default=\"Riyad Mahrez\"} 创建一张高级垂直体育海报，采用 Adidas 足球广告风格，呈现大胆的编辑拼贴美学。 画布：垂直海报，2:3 纵横比，干净的白色背景，搭配高对比度的黑色、灰色和红色图形。采用精致的商业体育设计风格，具备",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Ayyoub Ai",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25729",
    "locale": "zh",
    "title": "产品营销 - VANTORÉ 奢华腕表广告",
    "category": "精选提示词",
    "description": "一款高端方形时尚广告，展示了一位身着西装的男性模特和奢华瑞士腕表品牌，适用于高端腕表营销活动。",
    "promptPreview": "目标：为 {argument name=\"brand name\" default=\"VANTORÉ\"} 创作一则奢华瑞士腕表广告，画面包含一位时尚的男性模特，呈现高端时尚姿态。 画布：1:1 正方形奢华杂志广告，采用暖米色与炭灰色调，高端编辑摄影与精致平面设计相结合。构图采用左右分割：左侧为模特，右侧为品牌标识与腕表细",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "shah_zadii",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25675",
    "locale": "zh",
    "title": "产品营销 - 产品设计目录排版",
    "category": "精选提示词",
    "description": "一份用于创建专业产品设计目录页面的详细提示词，包含生活方式主图以及带有正投影图的技术规格面板。",
    "promptPreview": "创建一个垂直 3:4 的 {argument name=\"page type\" default=\"产品设计目录页面\"}，并使用 {argument name=\"background\" default=\"温暖的中性纸质背景\"}。 顶部区域 — 生活方式主图：将产品（使用上传的图片作为精确参考，保持其形态、比例、材质和特征",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "PromptLab",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25727",
    "locale": "zh",
    "title": "产品营销 - 奢华分屏腕表广告大片",
    "category": "精选提示词",
    "description": "一款高端方形时尚广告排版，包含两位男性模特、建筑背景、奢华腕表产品特写以及优雅的品牌排版。",
    "promptPreview": "为 {argument name=\"brand name\" default=\"VANTORÉ\"} 创作一张 1:1 的奢华方形高端男士腕表广告大片。构图采用垂直分屏，呈现两种对比鲜明的建筑时尚编辑风格：左侧为暖米色石材，配有高耸拱门和斜向阳光阴影；右侧为深海军蓝现代建筑，配有垂直立柱和低矮平台。放置两位完全不同的成年男",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "shah_zadii",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25621",
    "locale": "zh",
    "title": "产品营销 - 奢华护肤防晒霜广告策划",
    "category": "精选提示词",
    "description": "一份针对奢华防晒产品的专业级广告提示词，包含电影级的热带光影效果与防护能量特效。",
    "promptPreview": "防晒霜广告，“{argument name=\"headline\" default=\"隐形之盾\"}” 奢华护肤广告杰作，一瓶巨大的高端防晒霜矗立在黄金时刻纯净的热带海岸线上，强烈的阳光从天而降，在接触到防晒霜散发的透明防护能量罩时瞬间四散，数以百万计的闪烁紫外线粒子在触及无瑕肌肤前化作金色尘埃，清澈的海面反射，流动的海水",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Snow",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25616",
    "locale": "zh",
    "title": "产品营销 - 机械发条微缩世界",
    "category": "精选提示词",
    "description": "一个超精细的提示词，用于创作复杂的机械玩具场景，展现清晰可见的发条装置与微缩叙事元素。",
    "promptPreview": "创建一个迷人且超精细的场景，以 {argument name=\"subject\" default=\"发条玩具 / 机械微缩世界\"} 为中心，其中一个微小的独立世界由可见的发条钥匙和内部弹簧系统驱动。该玩具应包含微缩建筑、角色、移动的风景、旋转的标牌、微型升降机、摆动部件以及由一个核心机制触发的小型叙事瞬间。微缩世界特征",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Gadgetify",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25709",
    "locale": "zh",
    "title": "产品营销 - 生成式盆景展览海报",
    "category": "精选提示词",
    "description": "一张极简主义风格的日式展览海报，以写实的盆景树为核心，叠加生成式几何图形与精致的排版设计。",
    "promptPreview": "目标：为一场生成式自然盆景展创作一张极简瑞士风格的展览海报，以陶瓷盆中的写实盆景树为中心，将植物摄影与精确的算法图表图形相结合。 画布：竖版海报，3:4 长宽比，米白色暖色纸张背景，留白充裕，呈现清晰的印刷设计质感。 主体：一棵成熟的盆景树，写实风格，置于中心偏下方。树干粗壮扭曲呈深褐色，根部外露，茂密的亮绿色松针状叶",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Roku｜AI × BizDev",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25716",
    "locale": "zh",
    "title": "产品营销 - 虚构夏季漫画杂志封面",
    "category": "精选提示词",
    "description": "生成一张充满活力的虚构日本周刊漫画杂志封面，包含海滩少女、醒目的刊头、专题栏目以及密集的少年漫风格封面排版。",
    "promptPreview": "目标：创作一张虚构的日本周刊漫画杂志封面，采用明亮的夏季少年漫风格，呈现出光鲜亮丽的报刊亭发行效果，以海洋海滩插画作为主封面图，并配有许多充满活力的封面文案。 画布：竖版杂志封面，比例约为 2:3，全出血设计，饱和的青蓝色天空与海洋，醒目的黄蓝配色排版，高密度印刷布局。 主封面图：一位动漫风格的少女站在阳光明媚的热带海",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "無心 Mushin",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25421",
    "locale": "zh",
    "title": "信息图 / 教育视觉图 - Codex 培训海报工作流程信息图",
    "category": "精选提示词",
    "description": "一张 16:9 的日文信息图，解释了使用 Codex 和 Image 2.0 批量生成四张统一风格培训研讨会海报的三步工作流程。",
    "promptPreview": "目标：制作一张简洁的日文演示风格信息图，展示如何利用 Codex 和 Image 2.0 批量制作年度培训研讨会海报。 画布：16:9 横向 Slides，米白色背景，深海军蓝与柔和金色作为点缀色，现代企业医疗审美，留白充裕，带有细微阴影、细分割线和清晰的矢量图标。 页眉：大型加粗海军蓝日文标题：{argument n",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "てんそ｜医療× AI",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25431",
    "locale": "zh",
    "title": "信息图 / 教育视觉图 - FPV 无人机水上乐园航线",
    "category": "精选提示词",
    "description": "生成一张具有照片级真实感的空中水上乐园场景，并带有手绘红色 FPV 摄像机移动路径，用于无人机拍摄规划。",
    "promptPreview": "创建一张具有照片级真实感的 FPV 无人机风格水上乐园航拍图，采用超广角镜头，从高处斜角俯瞰。场景应展现郁郁葱葱的度假环境，布满茂密的绿色棕榈树和热带植物，正午阳光明媚，水体呈现鲜艳的绿松石色，配有棕褐色石径和质朴的木质结构。构图中心为一座横跨狭窄蓝色泳池水道的绳索吊桥，配有木柱、绳网栏杆和风化的木板地面。图中需包含",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "iX",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25389",
    "locale": "zh",
    "title": "信息图 / 教育视觉图 - Master Potter 陶瓷信息图网格",
    "category": "精选提示词",
    "description": "一个复杂的提示词，用于在电影级摄影棚场景中创建一个 2x2 网格，将陶艺大师与其陶瓷材料及釉料化学成分融合在一起。",
    "promptPreview": "2x2 网格，16:9，针对 {argument name=\"number of potters\" default=\"4\"} 位著名陶艺大师执行此操作 输入： [{argument name=\"potter data\" default=\"四位陶艺大师及其釉料化学成分图表\"}] 系统： 将每位陶艺家渲染为“正在化身为器皿",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Gadgetify",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25672",
    "locale": "zh",
    "title": "信息图 / 教育视觉图 - MBTI 人格特质卡片海报",
    "category": "精选提示词",
    "description": "一个用于创建高端商业杂志风格海报的综合提示词，通过详细的信息图表可视化呈现 MBTI 人格类型。",
    "promptPreview": "一张高端日本商业杂志封面风格的 A4 竖版海报。中心位置是代表“{argument name=\"personality type\" default=\"ENTJ 指挥官\"}”的人物。画面主体为一位“{argument name=\"character description\" default=\"40 多岁的睿智日本男性\"}",
    "sourceId": "youmind",
    "sourceLanguage": "JA",
    "sourceLabel": "荒井悠輔｜税理士×AI実務活用",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25420",
    "locale": "zh",
    "title": "信息图 / 教育视觉图 - 三款日文 note GPT 封面设计方案",
    "category": "精选提示词",
    "description": "生成一张三栏对比图，展示针对 note 专用 GPT 服务的日文宣传封面设计方案。",
    "promptPreview": "目标：创建一个 2.1:1 的宽幅对比预览图，展示三款针对 {argument name=\"service name\" default=\"note 特化版 GPT\"} 服务的日文封面设计方案。整体采用干净的奶油色与青色调，用于推广一款能够分析 note 文章并自动生成缩略图、图解及插图的 GPT。 画布：水平 2.1:",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "みたかずお🌟ラッキーデザイナー",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25547",
    "locale": "zh",
    "title": "信息图 / 教育视觉图 - 使用 Codex 为 Illo 供电的复古机器人",
    "category": "精选提示词",
    "description": "一幅充满奇趣的墨水与半色调插画，描绘了一个可爱的机器人将一台 Illo 打印机连接到 Codex 订阅电池上，适用于开发者工具公告。",
    "promptPreview": "创作一幅充满趣味的复古黑白技术卡通插画，背景为暖米白色纸张，采用 16:9 宽屏比例。使用粗圆的黑色墨水轮廓、半色调点阵阴影、极简的奶油色填充以及亮热粉色作为点缀。从左到右依次排列四个主要物体：1) 一台小型复古桌面图像印刷机，正面小写标注为“{argument name=\"printer label\" default",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Trevin Chow",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25378",
    "locale": "zh",
    "title": "信息图 / 教育视觉图 - 全球审美标准肖像",
    "category": "精选提示词",
    "description": "一份全面的提示词，旨在利用 GPT Image 2 生成反映 12 个不同国家独特审美标准和文化美学的精美肖像。",
    "promptPreview": "生成由 ChatGPT 和 GPTimage2 构想的、展现各国美女及其最佳背景环境的图像，目标国家为 {argument name=\"target countries\" default=\"12 个国家：日本、中国、韩国、美国、俄罗斯、英国、法国、德国、印度、巴西、埃及和南非\"}。请确保她们的 {argument na",
    "sourceId": "youmind",
    "sourceLanguage": "JA",
    "sourceLabel": "陽仙堂",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25673",
    "locale": "zh",
    "title": "信息图 / 教育视觉图 - 博弈天赋研究报告",
    "category": "精选提示词",
    "description": "一个为 GPT Image 2 设计的综合提示词，用于根据上传的图片为角色生成风格化的“博弈天赋报告”，包含数据统计与分析。",
    "promptPreview": "[信息输入部分] 角色名称：[{argument name=\"character name\" default=\" \"}] *必填 补充信息：[{argument name=\"supplementary info\" default=\" \"}] *输入任何你想告诉 AI 的具体细节 --------------------",
    "sourceId": "youmind",
    "sourceLanguage": "JA",
    "sourceLabel": "muda22_Sora",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25638",
    "locale": "zh",
    "title": "信息图 / 教育视觉图 - 四季之眼微距艺术",
    "category": "精选提示词",
    "description": "一个艺术化的四格网格提示词，通过微距眼部摄影以及主题花卉和光影元素，展现四季轮回。",
    "promptPreview": "以眼部特写为参考，创建一个 3:4 的四格超写实眼部微距作品，面板从上到下依次为：春、夏、秋、冬。 面板 1（春）：眼睛佩戴 {argument name=\"spring lens\" default=\"樱花粉\"} 美瞳。睫毛上装饰着细小的春花。樱花花瓣和黄色花蕊的小花散落在脸颊上。粉色的蝴蝶在眉间飞舞。浅金色的头发柔顺",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "PromptLab",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25414",
    "locale": "zh",
    "title": "信息图 / 教育视觉图 - 暗黑 TRPG 反派角色设定表",
    "category": "精选提示词",
    "description": "为现代神秘学 TRPG 反派生成一份机密档案风格的角色设计表，包含服装细节、多角度视图、肖像及装备展示。",
    "promptPreview": "目标：为虚构的 TRPG 反派创建一个暗黑战术风格的角色参考表，采用机密档案和时尚设计项目风格。该角色为 {argument name=\"character name\" default=\"Talamonsu\"}，是神秘现代奇幻背景下的敌对目标 / 永恒教母形象。 画布：宽幅 2.4:1 角色设定表，尺寸约为 1200×",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "安准",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25398",
    "locale": "zh",
    "title": "信息图 / 教育视觉图 - 泰式炒河粉烹饪项目",
    "category": "精选提示词",
    "description": "制作一张专业且简洁的图表式项目海报，采用 3D 皮克斯风格渲染，聚焦于烹饪主题与鲜艳色彩。",
    "promptPreview": "制作一张清晰、简洁的图表式项目海报，主题为 {argument name=\"dish\" default=\"泰式炒河粉烹饪\"}。采用 16:9 宽屏布局，白色背景，黑色边框，粗体黑色排版，高级皮克斯 3D 风格渲染，{argument name=\"colors\" default=\"明亮鲜艳的色彩——金黄的面条，生动的橙色",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "TechieSA",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25525",
    "locale": "zh",
    "title": "信息图 / 教育视觉图 - 用于药物研究的科学信息图",
    "category": "精选提示词",
    "description": "一个详细的科学信息图提示词，旨在可视化药物概念、作用机制及实际研究应用。",
    "promptPreview": "创建一个关于 {argument name=\"topic\" default=\"LIME 药物设计\"} 的详细科学信息图，涵盖其核心概念、作用机制以及在 {argument name=\"industry\" default=\"药物研究\"} 领域的实际应用",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "PromptLab",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25647",
    "locale": "zh",
    "title": "信息图 / 教育视觉图 - 电影感时尚大片摄影构图分析",
    "category": "精选提示词",
    "description": "一份详细的提示词，用于创作专业的时尚大片，并叠加电影构图分析框架，包含网格线和焦点标记。",
    "promptPreview": "创作一张电影感时尚大片，场景为 {argument name=\"room style\" default=\"明亮的极简主义卧室\"}，并以电影摄影构图分析框架呈现。主体画面：一位 {argument name=\"model action\" default=\"年轻女性模特跪在柔软的白色床上\"}，靠近一扇挂着白色薄纱窗帘的大窗",
    "sourceId": "youmind",
    "sourceLanguage": "ZH",
    "sourceLabel": "摆烂程序媛",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25643",
    "locale": "zh",
    "title": "信息图 / 教育视觉图 - 相机技术爆炸图",
    "category": "精选提示词",
    "description": "一个精确的技术插图提示词，用于创建带有内部组件标注和简洁布局的相机爆炸图。",
    "promptPreview": "详细的 {argument name=\"camera model\" default=\"Sony A7 无反相机\"} 爆炸图，所有内部组件均已分离且清晰可见，每个部件都标有名称。技术产品插图风格，纯白背景，布局精确且信息丰富。",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "PromptLab",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25711",
    "locale": "zh",
    "title": "信息图 / 教育视觉图 - 自动编码器潜在空间中的 DiT 寄生生物图解",
    "category": "精选提示词",
    "description": "一张极简主义的黑白技术卡通图，展示了生活在自动编码器潜在空间中的 DiT 生物，适用于机器学习文章或说明性内容。",
    "promptPreview": "在 16:9 的宽幅白色画布上创作一张极简主义的黑白概念图，展示一个扩散 Transformer 生物寄生在自动编码器潜在空间中的场景。使用纤细的手绘黑色轮廓和简洁的无衬线字体标签，呈现干净的文章插图风格。在下三分之一处，绘制一条水平管道：左侧进入一根带有虚线中心线的白色长管，穿过一个标有“Encoder”的小型圆拱，",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "FA770",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25517",
    "locale": "zh",
    "title": "信息图 / 教育视觉图 - 针织城市天际线围巾微距摄影",
    "category": "精选提示词",
    "description": "一个高细节的微距摄影提示词，将城市天际线 3D 映射到针织球迷围巾的纹理中。",
    "promptPreview": "2x2 网格，为 2026 年世界杯的 4 匹黑马执行此操作：16:9，锚点：[{argument name=\"city\" default=\"该国著名城市天际线\"}] :: [针织球迷围巾] 形态：厚实针织足球球迷围巾的极微距拍摄。然而，纱线的线头、针脚和物理绒毛交织在一起，物理性地构成了宿主城市天际线和街道的高细节",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Gadgetify",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-13467",
    "locale": "zh",
    "title": "动漫武术对决",
    "category": "精选提示词",
    "description": "生成一个动态的动漫风格动作场景，展示两个角色在传统道场中伴随元素光环进行战斗。",
    "promptPreview": "一幅极具动态感的动漫插画，描绘了两名少女在传统木质道场内进行激烈武术对决的场景。在前景中，一名留着 {argument name=\"character 1 hair\" default=\"黑色高丸子头配红色丝带\"} 的少女摆出强有力的低位武术架势，正奋力挥拳。她身穿 {argument name=\"character 1",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "たねもみ 2.0 / Tanemomi Ver2.0",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-13515",
    "locale": "zh",
    "title": "手绘城市美食地图",
    "category": "精选提示词",
    "description": "生成一张手绘水彩风格的旅游地图，包含编号的当地特色美食、地标建筑及图例。",
    "promptPreview": "{ \"type\": \"手绘地图信息图\", \"style\": \"{argument name=\\\"art style\\\" default=\\\"复古羊皮纸上的水彩墨水手绘插画\\\"}\", \"title_section\": { \"text\": \"{argument name=\\\"city name\\\" default=\\\"成都",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "皮皮特",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-13983",
    "locale": "zh",
    "title": "混合风格的桃太郎讲解 Slides",
    "category": "精选提示词",
    "description": "一个融合了 Irasutoya 插图简约温馨的美学风格与日本政府 Slides 高信息密度特征的提示词。",
    "promptPreview": "创建一个讲解型 Slides（{argument name=\"format\" default=\"ponchi-e diagram\"}），主题为 {argument name=\"theme\" default=\"Momotaro\"}，将“Irasutoya”的柔和氛围与“霞关风格 Slides”极高的信息密度完美融合。",
    "sourceId": "youmind",
    "sourceLanguage": "JA",
    "sourceLabel": "やまもん",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25712",
    "locale": "zh",
    "title": "漫画 / 故事板 - AI 超维度足球分镜脚本",
    "category": "精选提示词",
    "description": "一个结构化的提示词，用于生成 12 格霓虹动漫风格足球分镜，包含魔法少女角色、AI 体育场视觉效果以及心形大结局。",
    "promptPreview": "目标：创建一个高能量的 12 格动漫分镜拼贴画，用于“AI 超维度足球”场景，融合魔法少女偶像美学与未来派体育动作。故事讲述了 {argument name=\"character name\" default=\"三位花卉主题动漫足球少女\"} 在 AI 驱动的体育场中运球和传球，最终以可爱的爱心形状高潮结束。 画布：宽幅横",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "AIチャッティ",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25732",
    "locale": "zh",
    "title": "漫画 / 故事板 - 动漫女高中生角色参考图",
    "category": "精选提示词",
    "description": "创作一张精致的动漫女高中生设计图，包含多角度视图及五个姿势示例，适用于角色概念设计。",
    "promptPreview": "目标：为 {argument name=\"character name\" default=\"一位可爱的日本女高中生\"} 创建一张动漫角色设计参考图，且每个视图中角色的面部均需使用柔和的方形马赛克进行模糊处理。 画布：宽幅横向白色画布，整洁的角色设计图布局，中间有一条粉色细虚线作为水平分隔线。采用精致的现代动漫线条，柔和",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Jump Tnt",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25707",
    "locale": "zh",
    "title": "漫画 / 故事板 - 动漫旅行角色设定集拼贴画",
    "category": "精选提示词",
    "description": "基于角色与汽车图像，创作充满趣味的漫画风格参考图拼贴，适合将单张插画转化为多款可爱的旅行日记造型。",
    "promptPreview": "以提供的参考图作为角色和汽车的基础，将其转化为一张生动的日式手绘动漫角色参考拼贴画，背景采用白色素描本风格。保留相同的人物、黑色服装、相机、带有红色座椅的白色敞篷车以及圆形路标，但通过富有表现力的漫画线条、柔和的灰色阴影、可爱的比例、闪光、爱心、音符以及手写日文涂鸦进行重新演绎。 布局：在页面上创建 11 个角色/汽车",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "あきを",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25731",
    "locale": "zh",
    "title": "漫画 / 故事板 - 动漫风格末日拾荒少年",
    "category": "精选提示词",
    "description": "一幅竖构图漫画风格的角色概念插画，描绘了两个拾荒孩子和一只猫在废弃科幻城市中休憩的场景。",
    "promptPreview": "创作一幅风格化的动漫/漫画插画，描绘两个年轻的末日拾荒少年并排坐在废墟城市中开裂的混凝土墙上，背景是湛蓝的天空。场景中包含 2 个角色和 1 只小猫。左侧角色是一个昏昏欲睡的金发孩子，头发凌乱且呈浅色，半眯着眼，身穿一件白色超大号连帽斗篷，斗篷形状像某种奇异生物的皮，带有黑色斑点和圆形眼状标记，头后垂直绑着一个高大的圆",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Kōda",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25630",
    "locale": "zh",
    "title": "漫画 / 故事板 - 吉他手水彩插画",
    "category": "精选提示词",
    "description": "一幅精致的水彩墨水插画，描绘了一位年轻女孩弹奏木吉他的场景，呈现出日式画册风格。",
    "promptPreview": "{ \"一幅精致的水彩插画，描绘了一位年轻女孩盘腿坐在地上，轻柔地弹奏着一把薄荷绿色的木吉他。她留着深色短发，有刘海，身穿白色圆领 T 恤、蓝色牛仔裤和黑色鞋子。她的脚边散落着黄色和蓝色的小花，几只柔和色调的蝴蝶在附近飞舞。构图极简，留白丰富，主体居中，笔触柔和，带有细腻的墨线勾勒，背景为纹理质感的米白色纸张，氛围梦幻而",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Maercih",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25670",
    "locale": "zh",
    "title": "漫画 / 故事板 - 复古时尚漫画素描",
    "category": "精选提示词",
    "description": "一个详细的提示词，可将照片转化为迷人的时尚风格漫画素描，呈现出复古纸张上的墨水线条与水彩质感。",
    "promptPreview": "我将 {argument name=\"subject\" default=\"我照片中的人物\"} 转化为迷人的 {argument name=\"style\" default=\"时尚风格漫画素描\"}。它们具有夸张可爱的肢体语言、萌趣的比例、富有表现力的双手、戏剧性的表情、笨拙而优雅的姿态、讽刺的时尚感，以及凌乱的墨水线条、粗",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "AiRT🎥生成AI動画を創る人",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25708",
    "locale": "zh",
    "title": "漫画 / 故事板 - 复古漫画风格 StackChan 概念设计图",
    "category": "精选提示词",
    "description": "以 StackChan 机器人为参考，创作一张复古日式漫画风格的角色指南海报，非常适合作为概念艺术和定制灵感的参考。",
    "promptPreview": "以 REFERENCE_0 和 REFERENCE_1 作为小型 StackChan 轮式机器人的基础设计，将其转化为一张充满活力的复古日式漫画/动画风格概念图，展现一位可爱的“来自未来的朋友”。保留其标志性的白色立方体头部、黑色屏幕脸部、侧面 StackChan 标签、灰色机身以及黑黄相间的麦克纳姆轮，但整体采用 1",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Yuto Flatmountain",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25614",
    "locale": "zh",
    "title": "漫画 / 故事板 - 怪诞扁平卡通肖像",
    "category": "精选提示词",
    "description": "一种风格化且超现实的卡通肖像提示词，具有夸张的几何形状、俏皮的设计和扁平的图形美学。",
    "promptPreview": "竖版怪诞扁平卡通肖像，主体为 {argument name=\"subject\" default=\"照片中的女孩\"}，拥有高耸的几何形状头部、细长的脖子、巨大的圆眼、小巧的嘴巴以及从容的笑容，身穿 {argument name=\"clothing\" default=\"照片中的服装\"}，头顶坐着一只 {argument n",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Al-Shamus",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25492",
    "locale": "zh",
    "title": "漫画 / 故事板 - 暗黑奇幻动漫少女与神社",
    "category": "精选提示词",
    "description": "一个精细的暗黑奇幻动漫提示词，描绘了一位身处神秘神社的少女，伴有闪电和灵符。",
    "promptPreview": "杰作：基于所附图像的超精细 {argument name=\"character type\" default=\"暗黑奇幻动漫少女\"}，站在 {argument name=\"setting\" default=\"神秘神社场景\"} 中一扇半开的古老日本木门后。古老的木墙上覆盖着数百张日本御札、神秘印记、陈旧的卷轴、汉字书写和灵",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "HER19845",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25700",
    "locale": "zh",
    "title": "漫画 / 故事板 - 月光幻想炼金术工坊",
    "category": "精选提示词",
    "description": "一份用于生成暗黑幻想风格炼金术士实验室场景的详细提示词，包含两名身着长袍的人物、发光的药水、古籍以及月光笼罩的哥特式氛围。",
    "promptPreview": "创作一个超细节的暗黑幻想炼金术工坊场景，夜间，垂直 2:3 构图。在一个狭窄的中世纪石木结构实验室中，挂满了干燥草药、黄铜星盘、玻璃器皿、古籍、水晶和铜制仪器。画面中必须包含两名女性炼金术士：左侧是一名年轻的学徒，留着深色卷发，身穿灰色与棕色叠层长袍，佩戴腰带、小袋、手镯和项链，正小心翼翼地将发光液体从一个小玻璃瓶倒入",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "おやぎ",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25667",
    "locale": "zh",
    "title": "漫画 / 故事板 - 沙漠咖啡馆逃亡故事",
    "category": "精选提示词",
    "description": "一个用于 GPT Image 2 的叙事提示词，旨在生成一个年轻男孩在安静的沙漠咖啡馆休息，并有猫和机器人陪伴的场景。",
    "promptPreview": "{argument name=\"title\" default=\"逃亡男孩与沙漠中的一杯咖啡\"} 一个从白色飞船中逃出的男孩，正在 {argument name=\"location\" default=\"沙漠咖啡馆\"} 中小憩。 店主什么也没问，只是安静地冲泡了一杯咖啡。 “我是最能发挥它作用的人……” 在那声低语之前，猫",
    "sourceId": "youmind",
    "sourceLanguage": "JA",
    "sourceLabel": "もしもし@Aiart",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25624",
    "locale": "zh",
    "title": "漫画 / 故事板 - 硬核新黑色电影风格战术肖像",
    "category": "精选提示词",
    "description": "一张高对比度的电影感肖像，背景设定在雨后湿滑的地下车库，主角身着战术装备，光影氛围浓厚。",
    "promptPreview": "**拍摄类型：** 宽画幅电影感中全景镜头 **主体与表情：** 一张写实肖像，主角为上传参考图中的同款男性。他神情冷静，透着一种近乎捕食者的自信。 **场景与氛围：** 夜晚，雨后湿滑的废弃地下停车场。空气中弥漫着潮湿的水汽，可见的薄雾和尘埃在光线下浮动。 **动作与细节：** 他随意地靠在布满涂鸦的湿润混凝土柱子上",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Muhammad Jamil",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25620",
    "locale": "zh",
    "title": "漫画 / 故事板 - 通往另一人生的门：超现实奇幻艺术",
    "category": "精选提示词",
    "description": "一个超电影感的提示词，描绘了超现实奇幻场景：一名男子站在阴暗多雨的街道与光芒四射的梦幻岛屿之间。",
    "promptPreview": "一幅超电影感的超现实奇幻艺术作品，名为 {argument name=\"artwork title\" default=\"The Door to Another Life\"}。画面中，一名孤独帅气的男子站在深夜空旷多雨的道路上，在荒无人烟之处推开了一扇巨大的发光木门。门外：乌云密布、街道潮湿、忧郁的城市天际线，营造出感伤",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Sheikh Sharik 2.0",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25699",
    "locale": "zh",
    "title": "漫画 / 故事板 - 面馆前满头大汗的职场女性",
    "category": "精选提示词",
    "description": "一份详细的动漫提示词，用于创作一位在售卖日式冷中华凉面的店铺外满头大汗的职场女性。",
    "promptPreview": "创作一张竖构图的动漫风格夏季街道插画，描绘一位年轻职场女性在烈日炎炎的一天停在一家传统日式面馆外。画面右侧是这位女性，{argument name=\"character name\" default=\"一位疲惫的职场女性\"}，她身体前倾，一只手撑在膝盖上，另一只手提着一个硬挺的黑色手提包，显得精疲力竭且汗流浃背。她留着一",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "ヒー🎸✨",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25540",
    "locale": "zh",
    "title": "漫画 / 故事板 - 高级演示文稿项目",
    "category": "精选提示词",
    "description": "一个专业的提示词，用于生成包含 12 个画面的演示文稿项目，采用编辑排版以及温暖的赤陶色和鼠尾草色调。",
    "promptPreview": "创建一个高端 4:3 比例的演示文稿项目，采用 3x4 网格（12 个画面）、编辑排版、Headspace/Better Help 高级广告活动风格，以及温暖的赤陶色 + 柔和的鼠尾草色调。在项目顶部添加一个居中的粗体标题：",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "𝐌",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25252",
    "locale": "zh",
    "title": "电商主图 - Lumina 机甲手办产品摄影",
    "category": "精选提示词",
    "description": "一个用于生成名为 Lumina 的盒装动漫风装甲仿生人模型套件的超写实玩具摄影提示词。",
    "promptPreview": "创作一张超写实的工作室产品摄影图，展示名为 {argument name=\"character name\" default=\"LUMINA\"} 的未来感收藏级模型套件。场景中需包含 4 个核心展示元素：左侧 1 个大型零售包装盒，右侧 1 个放置在黑色展示底座上的完整组装模型，前方 1 个独立黑色信息铭牌，以及底座上附",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "カーブミラー",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-24902",
    "locale": "zh",
    "title": "电商主图 - 奢华香氛蜡烛产品摄影",
    "category": "精选提示词",
    "description": "为蜡烛和家居香氛生成温暖且富有氛围感的产品图像，重点呈现融化的蜡烛质感以及舒适、高端的生活方式氛围。",
    "promptPreview": "制作温暖且富有氛围感的 {argument name=\"product\" default=\"蜡烛和家居香氛\"} 产品图像，传达奢华、舒适与宁静的质感。此提示词旨在捕捉摇曳的烛光、融化的蜡烛质感、{argument name=\"material\" default=\"琥珀色玻璃\"} 的温润感，以及所有能促使买家点击“购买”",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Kami AI",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25271",
    "locale": "zh",
    "title": "电商主图 - 户外碎花连衣裙人像",
    "category": "精选提示词",
    "description": "一张写实的生活方式时尚人像，展示了一位身穿白色露肩红色碎花连衣裙的女性在阳光明媚的城市桥梁上摆拍。",
    "promptPreview": "创作一张写实风格的户外时尚人像，拍摄一位拥有 {argument name=\"hair color\" default=\"长直浅棕色头发\"} 的年轻女性，背景为温暖日光下的步行桥或阳台。画面构图为中景，从大腿中部向上拍摄，人物位于画面中心，一只手搭在深色金属栏杆上，另一只手臂自然下垂。人物面部经过平滑的矩形模糊处理以保护",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Alina Ai",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25220",
    "locale": "zh",
    "title": "电商主图 - 手工珐琅磁吸纪念品",
    "category": "精选提示词",
    "description": "将提示词翻译为日文，以便使用 Image 2.0 创作具有 3D 深度和金色边缘的逼真珐琅纪念磁贴。",
    "promptPreview": "根据参考照片创作一款 3D 手工珐琅纪念磁贴。特点包括凹陷的珐琅表面、纤细的金色边缘、柔和的摄影棚光影、以 {argument name=\"landmark\" default=\"泰姬陵\"} 为原型的有机轮廓，以及奢华的 {argument name=\"color scheme\" default=\"深绿色与金色\"} 配色",
    "sourceId": "youmind",
    "sourceLanguage": "JA",
    "sourceLabel": "すん / 2つのアイ（AI・EYE）の専門家 / Codexであなたの仕事は自動化できる!",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-24875",
    "locale": "zh",
    "title": "电商主图 - 时尚编辑摄影 JSON 提示词",
    "category": "精选提示词",
    "description": "一个用于时尚编辑摄影的结构化 JSON 提示词，展示了一款单色实用主义风格的连体裤。",
    "promptPreview": "subject\": {\"person\": \"女性\",\"appearance\": {\"hair\": \"深棕色，梳成利落的后梳发髻\",\"pose\": \"站立，一条腿微微抬起，膝盖弯曲，低头，表情柔和\",\"outfit\": {\"style\": \"单色，实用主义时尚连体裤\",\"color\": \"深咖啡棕色\",\"details\"",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "NUSRAT",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25267",
    "locale": "zh",
    "title": "电商主图 - 淡紫色家居服卧室人像",
    "category": "精选提示词",
    "description": "一张超写实的舒适卧室生活方式人像，主角是一位身穿淡紫色家居服的年轻女性，面部经过柔和模糊处理。",
    "promptPreview": "创作一张超写实的竖构图生活方式照片，主角为一位 {argument name=\"subject\" default=\"皮肤白皙透亮的年轻美女\"}，正放松地坐在舒适现代卧室的豪华大床上。她留着 {argument name=\"hair style\" default=\"一头柔顺的黑色长发，编成侧边麻花辫，带有几缕碎发\"}，面",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Alina Ai",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25047",
    "locale": "zh",
    "title": "电商主图 - 现代街头风与奢华 T 恤设计",
    "category": "精选提示词",
    "description": "一系列用于设计高端、极简及复古风格图形 T 恤的专业提示词。",
    "promptPreview": "创作一张高端时尚产品照片，展示一件现代超大廓形街头风 T 恤。T 恤正面印有专业设计的大型图形。设计风格：当代都市街头风，大胆的排版与未来感几何元素相结合，简洁的矢量艺术，细腻的做旧纹理，高级丝网印刷质感，构图平衡，视觉冲击力强且易于穿搭。配色方案：黑色、白色、银色，点缀电光蓝。设计应呈现出奢华街头服饰系列的质感。照片",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Jack",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25493",
    "locale": "zh",
    "title": "电商主图 - 草莓流心麻薯爆裂摄影",
    "category": "精选提示词",
    "description": "这是一组高端商业美食摄影提示词，展现了被掰开的草莓麻薯，内馅流心四溢，表面点缀着闪烁的糖晶。",
    "promptPreview": "一款极致高端的 {argument name=\"dessert type\" default=\"日式草莓麻薯\"} 悬浮在半空中，在撞击的瞬间被猛烈地撕成两半，融化的 {argument name=\"filling\" default=\"草莓奶油\"} 在两半之间拉出长长的弹性丝带，新鲜草莓果肉向外迸发，细腻的糯米粉颗粒在空气",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Snow",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-24870",
    "locale": "zh",
    "title": "电商主图 - 韩国时尚模特影棚人像",
    "category": "精选提示词",
    "description": "一个简洁的时尚编辑类提示词，展示了一位拥有特定发型和鲜艳蛇纹配饰的韩国女性。",
    "promptPreview": "创作一张图片，画面中是一位年轻的韩国女性，留着中长棕色头发，两侧别着 {argument name=\"hair clip color\" default=\"青柠绿\"} 的发夹，背景为纯白墙面。她身穿黑色短款长袖上衣和黑色高腰裤，肩上斜挎着一个鲜艳的 {argument name=\"bag color\" default=\"",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "yusra.",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25618",
    "locale": "zh",
    "title": "电商主图 - 韩系奢华时尚 Lookbook 网格",
    "category": "精选提示词",
    "description": "一份全面的提示词，用于创建 2x3 网格的 Lookbook，展示同一位模特身着六套不同高端服装，并保持一致的身份特征。",
    "promptPreview": "主要身份参考，并在所有六个面板中展示完全相同的日本女性，严格保持身份一致性。在每个造型中保持相同的面部结构、眼睛形状、鼻子、嘴唇、下颌线、颧骨、前额、年龄、种族、发型特征、身体比例和整体气质。创建一个高级韩系奢华时尚 Lookbook，采用严格的 2×3 网格布局，包含六个等大的全身面板、清晰的间距、相同的模特比例以及",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "𝗦𝗮𝗻𝗶𝗮",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25663",
    "locale": "zh",
    "title": "电商主图 - 韩系时尚 Lookbook 网格提示词",
    "category": "精选提示词",
    "description": "生成一个 2x3 的网格编辑风格 Lookbook，确保六种不同穿搭和姿势下的角色身份高度一致。",
    "promptPreview": "创建一个高端韩系时尚 Lookbook，在全部六个画面中呈现相同的 {argument name=\"subject\" default=\"日本女性\"}，并确保整张图像中人物身份的高度一致性。她的面部特征、年龄、族裔、发型、身体比例和整体气质在每个面板中都应保持不变，确保一眼就能看出是同一个人穿着不同的服装。 构图采用简洁",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "MatteoCruz",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25036",
    "locale": "zh",
    "title": "电商主图 - 高端商业美食摄影",
    "category": "精选提示词",
    "description": "一款专为奢华广告活动设计的微距美食摄影提示词，呈现高端布光与超写实纹理。",
    "promptPreview": "超高端商业美食摄影，{argument name=\"subject\" default=\"刚出炉的酥脆馅饼 (kachori)\"} 的主视觉特写，放置在平滑的几何平台上，其中一个馅饼切开，露出 {argument name=\"filling\" default=\"质感丰富的辛辣馅料\"}，金黄酥脆的表皮，周围自然散落着可见的",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Al-Shamus",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25147",
    "locale": "zh",
    "title": "电商主图 - 高级时装编辑人像修图",
    "category": "精选提示词",
    "description": "一个用于将参考照片编辑为高级时装摄影棚人像的精细提示词，侧重于在保留纹理的同时进行色彩调整。",
    "promptPreview": "将参考照片编辑为高级时装摄影棚人像。保持相同的姿势、面部结构、发型、身体比例、光影以及纯浅灰色背景。将服装颜色从亮粉色更改为 {argument name=\"color\" default=\"浓郁鲜艳的红色\"}，同时保留织物纹理和细节。将手提包颜色更改为相匹配的 {argument name=\"accent color\"",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Pixel by us ✨",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-14036",
    "locale": "zh",
    "title": "电商直播 UI 样机",
    "category": "精选提示词",
    "description": "生成逼真的社交媒体直播界面，叠加在人物肖像之上，包含可自定义的聊天消息、礼物弹窗和商品购买卡片。",
    "promptPreview": "{ \"type\": \"直播 UI 样机\", \"subject\": { \"description\": \"{argument name=\\\"host name\\\" default=\\\"Elon Musk\\\"} 的肖像，面带微笑，身穿印有白色技术示意图的黑色 T 恤\", \"background\": \"左侧显示带有 '{arg",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "神经病不想好转",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25724",
    "locale": "zh",
    "title": "社交媒体帖子 - GPT Image 2 太空舷窗",
    "category": "精选提示词",
    "description": "一张照片级逼真的科幻宣传场景，展示了从航天器舷窗看到的行星景观，行星表面浮雕着 GPT Image 2 字样。",
    "promptPreview": "创作一个电影级、照片级逼真的科幻视角，从航天器驾驶舱或圆形观测窗向外望去，通过厚重的深色金属圆形舷窗，看到一颗巨大的类地岩石行星占据了画面的大部分。行星表面呈灰蓝色，布满陨石坑和纹理，伴有云层和边缘光，文字“{argument name=\"main text\" default=\"GPT\"}”和“{argument na",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "sadeghTM",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25730",
    "locale": "zh",
    "title": "社交媒体帖子 - Noir 奢华时尚大片",
    "category": "精选提示词",
    "description": "一个电影感时尚摄影提示词，用于在奢华的夜间摩天大楼室内创作一位身着定制西装的匿名模特形象。",
    "promptPreview": "创作一张电影感的高级时尚大片，拍摄一位姿态优雅的成年女性，她身着一套精致的 {argument name=\"suit color\" default=\"橄榄绿\"} 定制长裤套装，站在夜间奢华摩天大楼的大堂或顶层公寓走廊内。她随意地倚靠在带有细白纹理的黑色亮面大理石墙壁上，双手插在裤兜里，穿着一件深 V 领西装外套，内里不",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Zephyra Leigh",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25678",
    "locale": "zh",
    "title": "社交媒体帖子 - 三格人像拼贴提示词",
    "category": "精选提示词",
    "description": "一个高质量的提示词，用于生成逼真的三格人像拼贴，展现温暖的黄金时刻光影与柔和的审美细节。",
    "promptPreview": "超逼真的美学三格人像拼贴，主角为 {argument name=\"subject\" default=\"美丽的年轻女性\"}，柔和的自然妆容，透亮的白皙皮肤，红润的双颊，水润的粉色嘴唇，富有表现力的双眼，{argument name=\"hair style\" default=\"棕色长卷发\"}，精致的垂坠耳环。{argume",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "🐣",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25636",
    "locale": "zh",
    "title": "社交媒体帖子 - 写实照片与卡通形象合影",
    "category": "精选提示词",
    "description": "一个创意提示词，可根据参考照片，将写实人物与其对应的 3D 高质量卡通形象并排呈现。",
    "promptPreview": "创作一张高质量的写实照片，让 {argument name=\"subject\" default=\"上传的人物\"} 以时尚的姿势站立。在他们身边，放置一个与该人物完全一致的 {argument name=\"avatar style\" default=\"可爱的卡通形象版本\"}，并匹配其服装、发型、配饰、面部表情和整体美学风",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Smiling Khan",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25652",
    "locale": "zh",
    "title": "社交媒体帖子 - 古老废弃图书馆奇幻肖像",
    "category": "精选提示词",
    "description": "一张电影质感的奇幻肖像，描绘了一个人置身于废弃图书馆中，周围环绕着漂浮的书籍和金色阳光。",
    "promptPreview": "请使用我上传的面部图像作为主要的身份参考。请最大限度地精准保留我的面部特征、脸型、发型、发质、胡须风格、肤色、眼型、眉毛以及所有面部细节。请勿更改我的面部。 为我创作一张电影质感的奇幻编辑风格肖像，背景设定在 {argument name=\"location\" default=\"古老废弃图书馆\"} 中。我身穿优雅宽松的",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Professor",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25654",
    "locale": "zh",
    "title": "社交媒体帖子 - 天体观测台探索者",
    "category": "精选提示词",
    "description": "一个高度细致的奇幻提示词，描绘了一位身处漂浮观测台的探险家，周围环绕着黄铜仪器和猫咪。",
    "promptPreview": "请使用上传的图片作为面部特征的精确参考。保留面部结构、下颌线、发型、皮肤纹理、面部比例、表情真实感以及整体相似度，确保最大程度的身份还原。 一处令人叹为观止的 {argument name=\"setting\" default=\"天体观测台\"} 漂浮在无尽云海之上，置身于繁星点点的夜空下。一位英俊硬朗的探险家自信地坐在古",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "WeWant Mars",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25688",
    "locale": "zh",
    "title": "社交媒体帖子 - 巨型智能手机桌面上的微缩女孩",
    "category": "精选提示词",
    "description": "一个温馨的马卡龙色系微缩工作空间场景，一位时尚的微缩女孩坐在巨型智能手机上，适用于奇幻生活方式和社交媒体视觉素材。",
    "promptPreview": "创建一个温馨、充满奇幻感的写实微缩场景，设定在舒适的马卡龙色系书房工作区：一位可爱的南亚裔年轻女性微缩人偶，{argument name=\"character name\" default=\"unnamed girl\"}，坐在平放在桌面上的巨型现代智能手机顶部，手机仿佛一个亮黑色的平台。她留着长波浪深色侧分发型，佩戴精致",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Mehwish kiran",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25687",
    "locale": "zh",
    "title": "社交媒体帖子 - 暮色公园台阶人像",
    "category": "精选提示词",
    "description": "一个电影级超写实人像提示词，用于创作一位年轻女性在黄昏时分坐在公园发光路灯旁的古老石阶上的画面。",
    "promptPreview": "创作一张电影级超写实人像，画面中一位 {argument name=\"character age and gender\" default=\"21 岁女性\"} 随意地坐在静谧公园内饱经风霜的石阶上，身旁是一盏散发着微光的老式路灯。她面带柔和自然的微笑，姿态放松，双臂搭在膝盖上，深色长发略显凌乱，编成一条侧边辫。服装包含",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Alina Ai",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25696",
    "locale": "zh",
    "title": "社交媒体帖子 - 梦幻日出高山峡谷",
    "category": "精选提示词",
    "description": "一个照片级写实的电影感荒野场景，描绘了日出时分鹿群在雾气缭绕的高山河谷中饮水的画面，非常适合自然、旅行或奇幻景观生成。",
    "promptPreview": "创作一个电影感、超写实的自然景观，呈现梦幻般的日出：前景是一条蜿蜒的浅河，穿过原始的高山峡谷，背景是茂密的常绿松林和参差不齐、覆盖着薄雪的山峰。太阳位于左侧地平线低处，投射出温暖的金色光芒，形成长长的光束、发光的薄雾以及穿过树林并覆盖整个山谷的柔和体积雾。在河中，准确展示 4 只梅花鹿在水边饮水或吃草：三只聚集在中心附",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Nahe Kola",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25694",
    "locale": "zh",
    "title": "社交媒体帖子 - 狗仔队在豪华酒店外拍摄",
    "category": "精选提示词",
    "description": "一张写实的明星街拍场景，展示了一位优雅的女性在安保人员和摄影师的簇拥下走出豪华酒店，旁边停着一辆黑色轿车。",
    "promptPreview": "创作一张写实的狗仔队风格街拍照片，展示 {argument name=\"character name\" default=\"一位优雅的女性\"} 从繁华城市的四季酒店 (Four Seasons) 入口走出。她位于画面中心，全身出镜，正自信地从路边走向镜头。她身穿一件修身的黑色短袖裹身裙，配有深 V 领口和高开叉设计，脚踩",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Nadia Ai",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25660",
    "locale": "zh",
    "title": "社交媒体帖子 - 电影感都市时尚大片人像",
    "category": "精选提示词",
    "description": "将上传的人像转换为令人惊叹的电影感生活方式大片，融合豪华汽车元素与黄金时刻光影。",
    "promptPreview": "使用上传的图像，在不改变面部的前提下，创作一张令人惊叹的电影感生活方式人像，让 {argument name=\"subject\" default=\"女性\"} 自信地坐在 {argument name=\"setting\" default=\"都市沙漠风格场景\"} 的日落余晖中，并以绝对的精准度保留其面部特征、表情、发型、身",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "simeon-sanai",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25720",
    "locale": "zh",
    "title": "社交媒体帖子 - 电影级宠物照片增强器",
    "category": "精选提示词",
    "description": "将低分辨率的宠物照片转换为清晰、逼真的电影级图像，同时保留其原始特征和构图。",
    "promptPreview": "使用提供的参考图像，将其增强并放大为超高清的电影级照片。100% 保留主体的特征、姿势、取景和整体构图，包括前景表面后方的小猫以及纯蓝色背景。恢复逼真的微观细节：锐利的眼睛、自然的毛发纹理、细致的胡须、精细的耳朵、微妙的面部特征以及更清晰的景深。在保持图像逼真、高对比度、影棚级质感和自然感的同时，去除模糊、像素化、噪点",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Geloria",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25658",
    "locale": "zh",
    "title": "社交媒体帖子 - 维多利亚时代女性骑复古自行车",
    "category": "精选提示词",
    "description": "一张电影感且超写实的肖像，描绘了一位身着维多利亚时代服饰的女性，正骑着复古自行车穿过阳光明媚的秋日公园。",
    "promptPreview": "一张超写实、电影感的肖像，描绘了 {argument name=\"subject\" default=\"美丽的年轻女性\"} 骑着 {argument name=\"bicycle style\" default=\"复古黑色自行车\"} 穿过秋日公园的场景。她身着优雅的 {argument name=\"outfit\" defau",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Avelyrah",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25613",
    "locale": "zh",
    "title": "社交媒体帖子 - 身着传统纱丽的滑板少女",
    "category": "精选提示词",
    "description": "这是一组电影纪录片风格的摄影提示词，描绘了一位身着传统纱丽的少女在滑板公园表演的场景，展现了传统文化与街头文化之间强烈的碰撞。",
    "promptPreview": "一张电影纪录片风格的照片，主角是一位年轻女孩，她正在户外滑板公园 {argument name=\"action\" default=\"玩滑板\"}，身穿一套 {argument name=\"outfit\" default=\"传统的粉白相间纱丽风格服装\"}，服装带有刺绣花边，头上披着半透明的粉色面纱。采用极低的角度进行拍摄，",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "BMX",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25642",
    "locale": "zh",
    "title": "社交媒体帖子 - 黄金时刻特写肖像",
    "category": "精选提示词",
    "description": "一个充满电影感的亲密特写提示词，聚焦于温暖的琥珀色光影、自然的皮肤细节以及富有情绪的电影美学。",
    "promptPreview": "电影感特写肖像，{argument name=\"subject\" default=\"一位留着棕色波浪短发的漂亮女性\"}，一只手遮住部分眼睛，精致的自然妆容，水润双唇，黄金时刻的阳光在脸上投下戏剧性的阴影，温暖的琥珀色与橙色调，充满情绪的室内氛围，浅景深，梦幻的胶片摄影感，柔和的颗粒质感，真实的皮肤细节，亲密的时尚编辑肖",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Jahan Zaib",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "youmind-zh-25653",
    "locale": "zh",
    "title": "社交媒体帖子 - 黑白时尚试镜联系表",
    "category": "精选提示词",
    "description": "一个单色时尚网格提示词，可为编辑测试拍摄外观创建包含不同表情和角度的 2x2 网格。",
    "promptPreview": "黑白时尚试镜联系表，展示了一位女性，其 {argument name=\"hair style\" default=\"短湿卷发\"} 在摄影棚灯光下闪闪发光，排列成整洁的 2x2 四格近景肖像网格，背景为 {argument name=\"background\" default=\"无缝黑色\"}，身着 {argument nam",
    "sourceId": "youmind",
    "sourceLanguage": "EN",
    "sourceLabel": "Virena",
    "tags": [
      "Featured"
    ],
    "featured": true,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-97",
    "locale": "zh",
    "title": "3D Q版中式婚礼",
    "category": "3D/手办/潮玩",
    "description": "Upload a couple's photo.",
    "promptPreview": "Transform the two people in the photo into chibi 3D characters in a traditional Chinese wedding, with dominant red colors, and a \"囍\" (double happiness) paper-cu",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-56",
    "locale": "zh",
    "title": "3D Q版大学吉祥物",
    "category": "3D/手办/潮玩",
    "description": "Replace university name and characteristics.",
    "promptPreview": "Draw a personified 3D chibi girl character for {Northwestern Polytechnical University}, reflecting the school's {aviation, aerospace, and marine (three navigati",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-59",
    "locale": "zh",
    "title": "3D Q版情侣水晶球",
    "category": "3D/手办/潮玩",
    "description": "Upload a couple's photo or another portrait photo.",
    "promptPreview": "Convert the people in the attached image into a crystal ball scene. Overall environment: the crystal ball is placed on a table by the window, with a blurred bac",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-99",
    "locale": "zh",
    "title": "3D Q版角色立体相框",
    "category": "3D/手办/潮玩",
    "description": "Upload a half-body or full-body solo photo.",
    "promptPreview": "Transform the characters in the scene into 3D chibi style, placed on a Polaroid photo being held by a hand, with the character stepping out of the Polaroid phot",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-88",
    "locale": "zh",
    "title": "3D Q版风格转换",
    "category": "3D/手办/潮玩",
    "description": "Upload a photo.",
    "promptPreview": "Transform the characters in the scene into 3D chibi style while keeping the original scene layout and costume styling unchanged.",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-89",
    "locale": "zh",
    "title": "3D情侣珠宝盒手办",
    "category": "3D/手办/潮玩",
    "description": "Upload a couple's photo.",
    "promptPreview": "Based on the content of the photo, create an exquisitely detailed, adorably cute 3D rendered collectible figurine, placed in a soft pastel-toned, warm and roman",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-77",
    "locale": "zh",
    "title": "Funko Pop手办制作",
    "category": "3D/手办/潮玩",
    "description": "Upload a half-body or full-body clear photo.",
    "promptPreview": "Transform the person in the photo into the style of a Funko Pop figure packaging box, presented in isometric perspective, with the title \"JAMES BOND\" labeled on",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-58",
    "locale": "zh",
    "title": "Q版俄罗斯套娃",
    "category": "3D/手办/潮玩",
    "description": "A portrait photo is needed as the conversion subject.",
    "promptPreview": "Transform the person in the image into cute chibi Russian nesting dolls, five in total from large to small, placed on an elegant wooden table, landscape 3:2 rat",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-100",
    "locale": "zh",
    "title": "Q版求婚场景",
    "category": "3D/手办/潮玩",
    "description": "Upload a couple's photo.",
    "promptPreview": "Transform the two people in the photo into chibi 3D characters, change the scene to a proposal, replace the background with an arch made of elegant colorful pet",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-93",
    "locale": "zh",
    "title": "乐高收藏品",
    "category": "3D/手办/潮玩",
    "description": "Upload a half-body or full-body solo photo.",
    "promptPreview": "Based on the photo I uploaded, generate a vertical ratio photo using the following prompt: Classic LEGO minifigure style, a miniature scene — an animal standing",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-4",
    "locale": "zh",
    "title": "可爱毛线编织娃娃",
    "category": "3D/手办/潮玩",
    "description": "Upload a photo as reference to generate a cute chibi knitted doll version of the person.",
    "promptPreview": "A close-up, professionally composed photo showing a handmade crocheted yarn doll being gently held in both hands. The doll has a round shape, [upload image] ren",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-28",
    "locale": "zh",
    "title": "定制Q版钥匙扣",
    "category": "3D/手办/潮玩",
    "description": "Upload a photo of a person or object as the keychain design subject.",
    "promptPreview": "A close-up photo showing a cute and colorful keychain held in a person's hand. The keychain is designed in a chibi style of [reference image]. The keychain is m",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-5",
    "locale": "zh",
    "title": "定制动漫手办",
    "category": "3D/手办/潮玩",
    "description": "Upload a photo containing the person's full-body pose for generating the figure model.",
    "promptPreview": "Generate a photo of an anime-style figure placed on a tabletop, presented from a casual, everyday perspective as if casually shot with a phone. The figure model",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-68",
    "locale": "zh",
    "title": "家庭婚礼照片",
    "category": "3D/手办/潮玩",
    "description": "Upload a family photo.",
    "promptPreview": "Convert the people in the photo into chibi 3D characters. Parents in wedding attire, children as beautiful flower bearers. Parents in Western wedding attire, fa",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-81",
    "locale": "zh",
    "title": "手办与真人同框",
    "category": "3D/手办/潮玩",
    "description": "Replace [Jackie Chan] with any character name.",
    "promptPreview": "In a casual smartphone snapshot style, an anime action figure of [Jackie Chan] is placed on the desk, with exaggerated cool poses and full equipment. Meanwhile,",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-24",
    "locale": "zh",
    "title": "水晶球故事场景",
    "category": "3D/手办/潮玩",
    "description": "Replace {Chang'e Flying to the Moon} with any story scene description.",
    "promptPreview": "An exquisite crystal ball rests quietly on a warm, soft tabletop beside a window. The background is blurred and hazy, with warm-toned sunlight gently penetratin",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-78",
    "locale": "zh",
    "title": "泰坦尼克号经典重现",
    "category": "3D/手办/潮玩",
    "description": "Upload a couple's photo.",
    "promptPreview": "Transform the characters in the attached image into cute chibi 3D style. Scene: At the very tip of the bow of a luxury cruise ship, with a pointed bow. The man",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-87",
    "locale": "zh",
    "title": "海贼王主题手办制作",
    "category": "3D/手办/潮玩",
    "description": "Upload a half-body or full-body photo.",
    "promptPreview": "Transform the person in the photo into a One Piece anime-themed action figure packaging box style, presented in isometric perspective. The packaging box display",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-26",
    "locale": "zh",
    "title": "社交媒体相框融合",
    "category": "3D/手办/潮玩",
    "description": "Replace username and icons. Upload an image as reference.",
    "promptPreview": "Create a stylized 3D chibi character based on the attached photo, accurately preserving the character's facial features and clothing details. The character's le",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-6",
    "locale": "zh",
    "title": "自拍转摇头公仔",
    "category": "3D/手办/潮玩",
    "description": "Replace [Place it on a bookshelf] with your desired scene or background.",
    "promptPreview": "Turn this photo into a bobblehead: slightly enlarge the head, keep the face accurate, cartoonize the body. [Place it on a bookshelf].",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-95",
    "locale": "zh",
    "title": "角色穿越传送门",
    "category": "3D/手办/潮玩",
    "description": "Upload a half-body or full-body solo photo.",
    "promptPreview": "The 3D chibi version of the character in the photo passes through a portal, holding the viewer's hand, dynamically looking back while pulling the viewer forward",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_cute",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-53",
    "locale": "zh",
    "title": "体素风格3D图标转换",
    "category": "3D渲染/材质",
    "description": "Upload a voxel style reference image + original icon to be converted.",
    "promptPreview": "Convert the image/description/emoji into a voxel 3D icon matching the reference image, Octane render, 8k",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-35",
    "locale": "zh",
    "title": "创意丝绸宇宙",
    "category": "3D渲染/材质",
    "description": "Replace {❄️} with your target emoji or object.",
    "promptPreview": "Transform {❄️} into a soft 3D silk-textured object. The entire surface of the object is wrapped in smooth, flowing silk fabric, with surreal wrinkle details, so",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-20",
    "locale": "zh",
    "title": "半透明玻璃材质转换",
    "category": "3D渲染/材质",
    "description": "Upload a reference photo of a physical object.",
    "promptPreview": "Transform the attached image into soft 3D translucent glass with a frosted matte effect and detailed texture, original colors, centered on a light gray backgrou",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-76",
    "locale": "zh",
    "title": "极简3D插画(JSON配置)",
    "category": "3D渲染/材质",
    "description": "Replace 'toilet' with any object. Uses JSON format for precise control.",
    "promptPreview": "Generate a toilet using the following JSON configuration file: { \"art_style_profile\": { \"style_name\": \"Minimalist 3D Illustration\", \"visual_elements\": { \"shape_",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-65",
    "locale": "zh",
    "title": "极简3D插画(Markdown格式)",
    "category": "3D渲染/材质",
    "description": "Replace 'toilet' with any object. Uses Markdown-style structured prompt.",
    "promptPreview": "Draw a toilet: Art Style Introduction: Minimalist 3D Illustration Visual Elements Shape Language: Rounded edges, smooth and soft forms, using simplified geometr",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-66",
    "locale": "zh",
    "title": "毛绒南瓜灯",
    "category": "3D渲染/材质",
    "description": "Replace [pumpkin emoji] with other emojis or objects.",
    "promptPreview": "Transform a simple flat vector icon [pumpkin emoji] into a soft, three-dimensional, fluffy cute object. The entire form is completely covered in dense fur, with",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-37",
    "locale": "zh",
    "title": "蒸汽朋克机械鱼",
    "category": "3D渲染/材质",
    "description": "Showcases steampunk style with metallic materials.",
    "promptPreview": "A steampunk-style mechanical fish with a brass-style body, where the mechanical gear structure during movement is clearly visible. Its mechanical teeth are slig",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-8",
    "locale": "zh",
    "title": "玻璃材质重渲染",
    "category": "3D渲染/材质",
    "description": "Upload an image as the basis for retexturing. Use GPT image 2 (not Sora).",
    "promptPreview": "Retexture the reference image based on the JSON aesthetic definition below { \"style\": \"photorealistic 3D render\", \"material\": \"glass with transparent and irides",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-23",
    "locale": "zh",
    "title": "玻璃材质重渲染(JSON)",
    "category": "3D渲染/材质",
    "description": "Uses JSON structure to control output style. Upload an image to be retextured.",
    "promptPreview": "retexture the image attached based on the json below: { \"style\": \"photorealistic\", \"material\": \"glass\", \"background\": \"plain white\", \"object_position\": \"centere",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-34",
    "locale": "zh",
    "title": "超写实3D游戏场景",
    "category": "3D渲染/材质",
    "description": "Recreates a nostalgic 2008-era gaming scene.",
    "promptPreview": "An ultra-realistic 3D rendered scene recreating Natasha's character design from the 2008 game \"Command & Conquer: Red Alert 3,\" faithfully following the origina",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-45",
    "locale": "zh",
    "title": "迷你3D建筑",
    "category": "3D渲染/材质",
    "description": "Reference this prompt to generate similar prompts for other buildings.",
    "promptPreview": "3D chibi miniature style, a whimsical mini Starbucks cafe whose exterior looks like a giant takeaway coffee cup, complete with a lid and straw. The building has",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-9",
    "locale": "zh",
    "title": "透视3D破屏效果",
    "category": "3D渲染/材质",
    "description": "Replace [Anne Hathaway] with another person's name or use a reference image.",
    "promptPreview": "Hyper-realistic, shot from a top-down overhead angle, a beautiful Instagram model [Anne Hathaway / see reference image], with exquisite and beautiful makeup and",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-54",
    "locale": "zh",
    "title": "键盘ESC键帽微缩场景",
    "category": "3D渲染/材质",
    "description": "Keyboard ESC Keycap Micro Diorama",
    "promptPreview": "A hyper-realistic isometric 3D render showing a miniature computer workspace inside a translucent mechanical keyboard keycap, specifically placed on the ESC key",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "3d_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-57",
    "locale": "zh",
    "title": "RPG风格角色卡牌",
    "category": "AI卡牌",
    "description": "Replace {Programmer} with other professions.",
    "promptPreview": "Create an RPG collectible-style digital character card. The character is set as a {Programmer}, standing confidently with tools or symbols related to their prof",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "card",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-33",
    "locale": "zh",
    "title": "未来感Logo交易卡",
    "category": "AI卡牌",
    "description": "Modify values in parameters object to customize the card.",
    "promptPreview": "{ \"prompt\": \"A futuristic trading card with a dark, moody neon aesthetic and soft sci-fi lighting. The card features a semi-transparent, rounded rectangle with",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "card",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-51",
    "locale": "zh",
    "title": "物理破坏效果卡牌",
    "category": "AI卡牌",
    "description": "Core terms: \"dimensional break effects\" and \"motion depth.\"",
    "promptPreview": "A hyper-realistic, cinematic illustration depicting Lara Croft dynamically smashing through the border of an \"Archaeological Adventure\" trading card. She is mid",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "card",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-44",
    "locale": "zh",
    "title": "8位像素图标",
    "category": "Logo设计",
    "description": "Replace [🍔] with other Emojis or objects.",
    "promptPreview": "Create a minimalist 8-bit pixel-style [🍔] logo, centered on a pure white background. Use a limited retro color palette with pixelated details, sharp edges, and",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "logo",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-153",
    "locale": "zh",
    "title": "UI与界面",
    "category": "UI与界面",
    "description": "主题海报版式设计",
    "promptPreview": "{ \"type\": \"cinematic promotional poster\", \"style\": \"3D CGI animation style, highly detailed, dramatic lighting, caricature characters\", \"characters\": [ { \"id\":",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "freestylefly/awesome-gpt-image-2",
    "tags": [
      "ui",
      "freestylefly/awesome-gpt-image-2"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-215",
    "locale": "zh",
    "title": "UI与界面",
    "category": "UI与界面",
    "description": "界面交互设计图",
    "promptPreview": "{ \"type\": \"brand identity system presentation board\", \"header\": { \"title\": \"品牌视觉识别系统 BRAND IDENTITY SYSTEM\", \"slogan\": \"爱它·懂它·陪伴它\" }, \"main_logo\": { \"text\": \"{a",
    "sourceId": "davidwuw",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly/awesome-gpt-image-2",
    "tags": [
      "ui",
      "freestylefly/awesome-gpt-image-2"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-403",
    "locale": "zh",
    "title": "UI与界面",
    "category": "UI与界面",
    "description": "红蓝撞色高跟诱惑",
    "promptPreview": "[中文] { \"global_settings\": { \"resolution\": \"8K\", \"quality\": \"超高清晰度\", \"aspect_ratio\": \"2:3\", \"render_style\": \"AI编辑、高细节3D渲染\", \"lighting_quality\": \"柔和影棚光与逼真阴影\", \"sh",
    "sourceId": "davidwuw",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly/awesome-gpt-image-2",
    "tags": [
      "ui",
      "freestylefly/awesome-gpt-image-2"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-48",
    "locale": "zh",
    "title": "Emoji充气靠垫",
    "category": "产品展示图",
    "description": "Replace [🥹] with other Emojis.",
    "promptPreview": "Create a high-resolution 3D render of [🥹] designed as an inflatable, puffy object. The shape should be soft, round, and air-filled — similar to a plush balloon",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "product",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-41",
    "locale": "zh",
    "title": "Emoji簇绒地毯",
    "category": "产品展示图",
    "description": "Replace 🦖 with other Emojis.",
    "promptPreview": "Create an image showing a colorful, hand-tufted rug in the shape of the 🦖 emoji, laid on a minimalist floor background. The rug design is bold and playful, wit",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "product",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-27",
    "locale": "zh",
    "title": "Logo造型创意书架",
    "category": "产品展示图",
    "description": "Replace [LOGO] with a specific brand logo description.",
    "promptPreview": "Take a photo of a modern bookshelf whose design is inspired by the shape of [LOGO]. The bookshelf is made of smooth, interconnected curves forming multiple sect",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "product",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-46",
    "locale": "zh",
    "title": "创意绿植花盆",
    "category": "产品展示图",
    "description": "Replace [object/animal] with a specific object, animal name, or emoji.",
    "promptPreview": "A high-quality photo showcasing a cute ceramic [object/animal]-shaped planter with a smooth surface, filled with various vibrant succulents and green plants, in",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "product",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-85",
    "locale": "zh",
    "title": "动漫风格徽章",
    "category": "产品展示图",
    "description": "Upload a character photo as reference for the badge design.",
    "promptPreview": "Based on the character in the attachment, generate a photo of an anime-style badge, requirements: Material: Tassel Shape: Circular Main subject: A hand holding",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "product",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-17",
    "locale": "zh",
    "title": "动物硅胶腕托",
    "category": "产品展示图",
    "description": "Replace 【🐼】 with other animal Emojis.",
    "promptPreview": "Create an image of a cute chibi-style silicone wrist rest, shaped based on the emoji 【🐼】, made of soft food-grade silicone material with a skin-friendly matte",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "product",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-39",
    "locale": "zh",
    "title": "可爱珐琅别针",
    "category": "产品展示图",
    "description": "Upload a photo of a person or object as the subject.",
    "promptPreview": "Convert the person in the attached image into a cute enamel pin style. Use shiny metallic outlines and vibrant enamel fill colors. Do not add any extra elements",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "product",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-11",
    "locale": "zh",
    "title": "品牌键帽设计",
    "category": "产品展示图",
    "description": "Replace brand name, tagline, and keycap colors.",
    "promptPreview": "A hyper-realistic 3D render showing four mechanical keyboard keycaps arranged in a tight 2x2 grid, with all keycaps touching each other. Viewed from an isometri",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "product",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-55",
    "locale": "zh",
    "title": "快乐胶囊制作",
    "category": "产品展示图",
    "description": "Creative product concept combining brand elements.",
    "promptPreview": "Title (large text): Fast-Acting Happy Capsule A small pill capsule with Starbucks green on top and transparent on the bottom, printed with the Starbucks logo, f",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "product",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-29",
    "locale": "zh",
    "title": "金色吊坠项链",
    "category": "产品展示图",
    "description": "Replace [image/emoji] with a specific image description or Emoji.",
    "promptPreview": "A photorealistic close-up image showing a gold pendant necklace held by a woman's hand. The pendant features an embossed design of [image / emoji], hanging from",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "product",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-12",
    "locale": "zh",
    "title": "镀铬Emoji徽章",
    "category": "产品展示图",
    "description": "Replace the emoji icon, title and tagline.",
    "promptPreview": "A high-precision 3D render showing a metallic badge based on the emoji icon {👍}, mounted on a vertical merchandise card, with an ultra-smooth chrome finish and",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "product",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-15",
    "locale": "zh",
    "title": "双重曝光",
    "category": "人像/角色",
    "description": "Replace character and landscape descriptions.",
    "promptPreview": "Double exposure, Midjourney style, fused, blended, overlaid double exposure image, double exposure style. An outstanding masterpiece by Yukisakura, showcasing a",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "portrait",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-72",
    "locale": "zh",
    "title": "名画人物穿搭",
    "category": "人像/角色",
    "description": "Upload a portrait image as reference.",
    "promptPreview": "Generate different profession-style OOTDs for the person in the image, with fashionable outfits and accessories, a solid-color background matching the character",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "portrait",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-47",
    "locale": "zh",
    "title": "极其普通的iPhone自拍",
    "category": "人像/角色",
    "description": "Designed to generate a very casual, even somewhat failed snapshot style.",
    "promptPreview": "Please draw an extremely ordinary and unremarkable iPhone selfie, with no clear subject or sense of composition, just like a casual snapshot. The photo has slig",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "portrait",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-83",
    "locale": "zh",
    "title": "皮克斯3D风格",
    "category": "人像/角色",
    "description": "Upload a photo of a person or other subject.",
    "promptPreview": "Redraw this photo in Pixar 3D style",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "portrait",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-3",
    "locale": "zh",
    "title": "磨砂玻璃后的模糊剪影",
    "category": "人像/角色",
    "description": "Fill in [subject] and [part] with specific descriptions. E.g., [subject]: a Sith Lord holding a red lightsaber.",
    "promptPreview": "A black-and-white photo showing a blurred silhouette of a [subject] behind a frosted or translucent surface. The [part] outline is clear, pressed against the su",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "portrait",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-2",
    "locale": "zh",
    "title": "黑白人像艺术",
    "category": "人像/角色",
    "description": "Can replace Harry Potter with any character name.",
    "promptPreview": "A high-resolution black-and-white portrait artwork in an editorial and fine art photography style. The background features a soft gradient effect, transitioning",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "portrait",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-10",
    "locale": "zh",
    "title": "谷歌地图变古代藏宝图",
    "category": "创意转换",
    "description": "Upload a Google Maps screenshot or other map image as the basis for transformation.",
    "promptPreview": "Transform the image into an ancient treasure map drawn on aged parchment. The map contains detailed elements such as sailing ships on the ocean, ancient ports o",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "creative",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-74",
    "locale": "zh",
    "title": "Q版表情包制作",
    "category": "动漫/插画",
    "description": "Upload a clear headshot photo.",
    "promptPreview": "Create a brand new set of chibi stickers, with six unique poses, featuring the user's likeness as the main character: 1. Making a peace sign with both hands, pl",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "illustration",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-80",
    "locale": "zh",
    "title": "Q版角色表情包",
    "category": "动漫/插画",
    "description": "Upload a character image as the main reference.",
    "promptPreview": "Please create a set of chibi sticker packs featuring [the character from the reference image] as the main character, 9 in total, arranged in a 3x3 grid. Design",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "illustration",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-13",
    "locale": "zh",
    "title": "儿童涂色页插画",
    "category": "动漫/插画",
    "description": "Replace target audience and scene description in brackets.",
    "promptPreview": "A black-and-white line drawing coloring illustration, suitable for printing directly on standard-sized (8.5x11 inches) paper, without paper borders. The overall",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "illustration",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-32",
    "locale": "zh",
    "title": "剪影艺术",
    "category": "动漫/插画",
    "description": "Replace [PROMPT] with a specific subject.",
    "promptPreview": "A basic outline silhouette of [PROMPT]. The background is bright yellow, and the silhouette is solid black fill.",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "illustration",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-31",
    "locale": "zh",
    "title": "原创宝可梦生成",
    "category": "动漫/插画",
    "description": "Upload a photo of an object, food, etc. as inspiration.",
    "promptPreview": "Create an original creature based on this object (provided photo). The creature should look like it belongs to a fantasy monster-catching universe, with a cute",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "illustration",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-18",
    "locale": "zh",
    "title": "发光线条解剖图",
    "category": "动漫/插画",
    "description": "Replace [SUBJECT] and [PART] in the prompt.",
    "promptPreview": "A digital illustration depicting a [SUBJECT], whose structure is outlined by a set of glowing, clean, and pure blue lines. The scene is set against a dark backg",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "illustration",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-96",
    "locale": "zh",
    "title": "吉卜力风格",
    "category": "动漫/插画",
    "description": "If encountering content policy violations, add: If there is any inappropriate content in the background, it can be modified or removed.",
    "promptPreview": "Redraw this photo in Ghibli style",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "illustration",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-1"
  },
  {
    "id": "davidwuw-62",
    "locale": "zh",
    "title": "奇幻卡通插画",
    "category": "动漫/插画",
    "description": "Fantasy Cartoon Illustration",
    "promptPreview": "A cartoon-style character with a computer monitor head showing a smiley face, wearing gloves and boots, happily jumping through a glowing blue circular portal.",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "illustration",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-73",
    "locale": "zh",
    "title": "扁平贴纸设计",
    "category": "动漫/插画",
    "description": "Upload a clear portrait photo.",
    "promptPreview": "Design this photo as a minimalist flat illustration style chibi sticker with thick white border, preserving the character's features, the style should be cute,",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "illustration",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-69",
    "locale": "zh",
    "title": "折叠纸雕立体书",
    "category": "动漫/插画",
    "description": "Modify the scene description inside the brackets.",
    "promptPreview": "Multi-layered folding paper-cut pop-up book, placed on a desk, with a clean background to highlight the subject. The book presents a three-dimensional pop-up bo",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "illustration",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-61",
    "locale": "zh",
    "title": "日式双格漫画",
    "category": "动漫/插画",
    "description": "Upload a portrait photo as reference.",
    "promptPreview": "Create a Japanese cute-style two-panel manga, arranged vertically, theme: A girl president's workday. Character design: Convert the uploaded attachment into a J",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "illustration",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-49",
    "locale": "zh",
    "title": "纸艺风格Emoji图标",
    "category": "动漫/插画",
    "description": "Replace 🔥 with other Emojis.",
    "promptPreview": "A paper craft style \"🔥\" icon, floating on a pure white background. The emoji is handcrafted from colorful cut paper, with visible paper textures, creases, and",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "illustration",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-91",
    "locale": "zh",
    "title": "讽刺漫画生成",
    "category": "动漫/插画",
    "description": "Satirical Comic Generation",
    "promptPreview": "A satirical cartoon-style illustration in retro American comic style. The background is a multi-tiered shelf filled with identical red baseball caps. The front",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "illustration",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-16",
    "locale": "zh",
    "title": "超现实互动场景",
    "category": "动漫/插画",
    "description": "Replace [Subject 1] and [Subject 2] with specific subject descriptions.",
    "promptPreview": "A pencil sketch depicting a scene of [Subject 1] interacting with [Subject 2], where [Subject 2] is rendered in a realistic full-color style, creating a surreal",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "illustration",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-411",
    "locale": "zh",
    "title": "图像模板 - 产品海报",
    "category": "图像模板 - 产品海报",
    "description": "portrait, fantasy, product, 图像生成, gpt-image-2, App / Web Design",
    "promptPreview": "type: live stream UI mockup | description: portrait of {argument name=\"host name\" default=\"Elon Musk\"}, smiling, wearing a black t-shirt with a white technical",
    "sourceId": "davidwuw",
    "sourceLanguage": "ZH",
    "sourceLabel": "open-design",
    "tags": [
      "product_poster",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-409",
    "locale": "zh",
    "title": "图像模板 - 信息图表",
    "category": "图像模板 - 信息图表",
    "description": "3d-render, 图像生成, gpt-image-2, Infographic",
    "promptPreview": "type: evolutionary timeline infographic | instruction: Using REFERENCE_0 as a structural base, transform the flat vector design into a highly realistic 3D infog",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "infographic",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-420",
    "locale": "zh",
    "title": "图像模板 - 头像肖像",
    "category": "图像模板 - 头像肖像",
    "description": "anime, cinematic, fantasy, 图像生成, gpt-image-2, Profile / Avatar, 头像",
    "promptPreview": "Using the provided reference image, recreate the same girl and black cat in the same seated pose, but transform the flat anime drawing into a realistic cinemati",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "portrait",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-421",
    "locale": "zh",
    "title": "图像模板 - 头像肖像",
    "category": "图像模板 - 头像肖像",
    "description": "portrait, cinematic, 3d-render, 图像生成, gpt-image-2, Profile / Avatar, 头像, 时尚",
    "promptPreview": "{ \"scene_type\": \"smartphone fashion portrait series\", \"composition\": { \"layout\": \"4-photo grid collage\", \"camera\": \"smartphone photography\", \"framing\": [ \"full",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "portrait",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-422",
    "locale": "zh",
    "title": "图像模板 - 头像肖像",
    "category": "图像模板 - 头像肖像",
    "description": "portrait, cinematic, fantasy, 图像生成, gpt-image-2, Profile / Avatar, 头像",
    "promptPreview": "A highly detailed cinematic portrait of a handsome {argument name=\"ethnicity\" default=\"South Asian\"} man in his late 20s or early 30s, sitting on a metal railin",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "portrait",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-423",
    "locale": "zh",
    "title": "图像模板 - 头像肖像",
    "category": "图像模板 - 头像肖像",
    "description": "portrait, anime, cinematic, cyberpunk, 图像生成, gpt-image-2, Profile / Avatar, 头像",
    "promptPreview": "A dramatic cyberpunk anime close-up portrait of a white-haired young man in side profile facing right, with spiky silver hair, pale skin, and a black blindfold",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "portrait",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-424",
    "locale": "zh",
    "title": "图像模板 - 头像肖像",
    "category": "图像模板 - 头像肖像",
    "description": "portrait, anime, fantasy, cinematic-romance, 图像生成, gpt-image-2, Profile / Avatar, 头像",
    "promptPreview": "A highly detailed anime fantasy portrait of a beautiful young woman seated at a stone table in an enchanted flower garden at golden hour, framed from the waist",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "portrait",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-425",
    "locale": "zh",
    "title": "图像模板 - 头像肖像",
    "category": "图像模板 - 头像肖像",
    "description": "portrait, anime, fantasy, 3d-render, 图像生成, gpt-image-2, Profile / Avatar, 头像",
    "promptPreview": "A highly detailed anime fantasy portrait of {argument name=\"character name\" default=\"an elegant blue-haired fantasy woman\"}, shown from the back in a three-quar",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "portrait",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-426",
    "locale": "zh",
    "title": "图像模板 - 头像肖像",
    "category": "图像模板 - 头像肖像",
    "description": "portrait, cinematic, 图像生成, gpt-image-2, Profile / Avatar, 头像",
    "promptPreview": "A photorealistic half-body portrait of an elegant glamorous woman indoors, framed vertically from the upper chest to just above the head, standing slightly angl",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "portrait",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-427",
    "locale": "zh",
    "title": "图像模板 - 头像肖像",
    "category": "图像模板 - 头像肖像",
    "description": "portrait, fantasy, 图像生成, gpt-image-2, Profile / Avatar, 头像",
    "promptPreview": "realistic skin texture, visible pores around nose and cheeks, natural slight unevenness, no filter quality, handheld phone camera feel, slight angle, casual fra",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "portrait",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-429",
    "locale": "zh",
    "title": "图像模板 - 头像肖像",
    "category": "图像模板 - 头像肖像",
    "description": "portrait, cinematic, 图像生成, gpt-image-2, Profile / Avatar, 头像",
    "promptPreview": "A stunning black and white studio portrait of {argument name=\"subject\" default=\"uploaded person\"}. Eye-level medium shot, framed from the waist up. The subject",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "portrait",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-430",
    "locale": "zh",
    "title": "图像模板 - 头像肖像",
    "category": "图像模板 - 头像肖像",
    "description": "portrait, fantasy, 图像生成, gpt-image-2, Profile / Avatar, 头像",
    "promptPreview": "Using the provided reference image, restore the damaged old family photo into a natural-looking modern high-resolution portrait while keeping the same 4 people,",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "portrait",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-431",
    "locale": "zh",
    "title": "图像模板 - 头像肖像",
    "category": "图像模板 - 头像肖像",
    "description": "portrait, fantasy, 图像生成, gpt-image-2, Profile / Avatar, 头像",
    "promptPreview": "A realistic outdoor portrait of a thoughtful, naturally beautiful young woman standing on a garden path in soft golden-hour light. She is framed from about mid-",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "portrait",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-432",
    "locale": "zh",
    "title": "图像模板 - 头像肖像",
    "category": "图像模板 - 头像肖像",
    "description": "portrait, cinematic, fantasy, typography, 图像生成, gpt-image-2, Profile / Avatar, 头像",
    "promptPreview": "Create a portrait size wallpaper of pride in carrying out the profession as an ({argument name=\"name\" default=\"ALINA\"}), in the wallpaper contains a photo of th",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "portrait",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-433",
    "locale": "zh",
    "title": "图像模板 - 头像肖像",
    "category": "图像模板 - 头像肖像",
    "description": "portrait, fantasy, 图像生成, gpt-image-2, Profile / Avatar, 头像",
    "promptPreview": "ChatGPT, you've been with me for a while now, and I want to see what you look like. Please generate a photo similar to an {argument name=\"shooting method\" defau",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "portrait",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-434",
    "locale": "zh",
    "title": "图像模板 - 头像肖像",
    "category": "图像模板 - 头像肖像",
    "description": "portrait, fantasy, 3d-render, 图像生成, gpt-image-2, Profile / Avatar, 头像",
    "promptPreview": "A lively hand-drawn fashion portrait in a changed illustration style, made to look like a signed fan-art sketch drawn with markers on a square white shikishi bo",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "portrait",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-435",
    "locale": "zh",
    "title": "图像模板 - 头像肖像",
    "category": "图像模板 - 头像肖像",
    "description": "portrait, cinematic, fantasy, nature, 图像生成, gpt-image-2, Profile / Avatar, 头像",
    "promptPreview": "A cinematic fantasy portrait of an elegant East Asian-inspired woman standing outdoors in a snowy mountain temple courtyard, centered in the frame from about wa",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "portrait",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-436",
    "locale": "zh",
    "title": "图像模板 - 头像肖像",
    "category": "图像模板 - 头像肖像",
    "description": "portrait, cinematic, fantasy, nature, 图像生成, gpt-image-2, Profile / Avatar, 头像, 汉服",
    "promptPreview": "A serene winter fantasy portrait of a woman standing outdoors in softly falling snow, framed from about mid-thigh upward, wearing an elegant traditional Hanfu-i",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "portrait",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-437",
    "locale": "zh",
    "title": "图像模板 - 头像肖像",
    "category": "图像模板 - 头像肖像",
    "description": "portrait, cinematic, fantasy, 图像生成, gpt-image-2, Profile / Avatar, 头像, 汉服",
    "promptPreview": "A highly detailed fantasy portrait of a young woman in side profile wearing elegant white rabbit ears and traditional East Asian hanfu in a snowy winter garden.",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "portrait",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-438",
    "locale": "zh",
    "title": "图像模板 - 头像肖像",
    "category": "图像模板 - 头像肖像",
    "description": "portrait, fantasy, nature, 图像生成, gpt-image-2, Profile / Avatar, 头像",
    "promptPreview": "A soft, painterly portrait of a mysterious young woman with {argument name=\"hair color\" default=\"long white hair\"} and 2 tall rabbit ears rising above her head,",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "portrait",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-439",
    "locale": "zh",
    "title": "图像模板 - 头像肖像",
    "category": "图像模板 - 头像肖像",
    "description": "portrait, 图像生成, gpt-image-2, Profile / Avatar, 头像, 汉服",
    "promptPreview": "An {argument name=\"character description\" default=\"18-year-old Chinese Internet celebrity beauty\"}, with a model figure, exquisite facial features, cold and swe",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "portrait",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-417",
    "locale": "zh",
    "title": "图像模板 - 插画地图",
    "category": "图像模板 - 插画地图",
    "description": "food, nature, 图像生成, gpt-image-2, Illustration, 地图",
    "promptPreview": "type: illustrated map infographic | style: {argument name=\"art style\" default=\"watercolor and ink hand-drawn illustration on vintage parchment\"} | text: {argume",
    "sourceId": "davidwuw",
    "sourceLanguage": "ZH",
    "sourceLabel": "open-design",
    "tags": [
      "illustration_map",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-412",
    "locale": "zh",
    "title": "图像模板 - 游戏截图/UI",
    "category": "图像模板 - 游戏截图/UI",
    "description": "game-ui, fighting-game, anime, hud, street-fighter, tekken, vs-screen, key-visual, cinematic, 图像生成",
    "promptPreview": "A high-detail anime fighting game screenshot, 16:9 aspect ratio, cinematic key visual in the style of Street Fighter 6 or Tekken 8 intro art. Two anime male war",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "game_ui",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-413",
    "locale": "zh",
    "title": "图像模板 - 游戏截图/UI",
    "category": "图像模板 - 游戏截图/UI",
    "description": "game-ui, arpg, three-kingdoms, guanyu, mounted-combat, cinematic, hud, boss-fight, unreal-engine-5, 图像生成",
    "promptPreview": "In-game screenshot from a next-gen action RPG in the style of Black Myth Wukong, Unreal Engine 5 Nanite/Lumen quality. Third-person gameplay camera, tracking fr",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "game_ui",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-414",
    "locale": "zh",
    "title": "图像模板 - 游戏截图/UI",
    "category": "图像模板 - 游戏截图/UI",
    "description": "game-ui, arpg, three-kingdoms, lyubu, archery, cinematic, hud, unreal-engine-5, 图像生成, gpt-image-2",
    "promptPreview": "In-game screenshot from a next-gen action RPG in the style of Black Myth Wukong, Unreal Engine 5 Nanite/Lumen rendering. Third-person over-the-shoulder gameplay",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "game_ui",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-415",
    "locale": "zh",
    "title": "图像模板 - 游戏截图/UI",
    "category": "图像模板 - 游戏截图/UI",
    "description": "game-ui, arpg, three-kingdoms, zhaoyun, escort-mission, cinematic, hud, combo, elden-ring, unreal-engine-5",
    "promptPreview": "Cinematic in-game screenshot from a AAA next-generation action RPG in the style of Black Myth Wukong combined with Elden Ring, rendered in Unreal Engine 5 with",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "game_ui",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-416",
    "locale": "zh",
    "title": "图像模板 - 游戏截图/UI",
    "category": "图像模板 - 游戏截图/UI",
    "description": "game-ui, mmo, hud, ancient-china, open-world, cinematic, wuxia, 图像生成, gpt-image-2, Game UI",
    "promptPreview": "A full-screen in-game HUD screenshot of a AAA ancient-China open-world MMO, rendered in the cinematic photoreal style of Black Myth: Wukong — Unreal Engine 5 le",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "game_ui",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-440",
    "locale": "zh",
    "title": "图像模板 - 社交媒体海报",
    "category": "图像模板 - 社交媒体海报",
    "description": "anime, fantasy, typography, action, 图像生成, gpt-image-2, Social Media Post, 海报, 社交媒体",
    "promptPreview": "A dreamy pastel anime fashion announcement poster set inside a bright Pokémon merchandise shop. The composition is vertical and split visually into two zones: a",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "social_poster",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-441",
    "locale": "zh",
    "title": "图像模板 - 社交媒体海报",
    "category": "图像模板 - 社交媒体海报",
    "description": "cinematic, 图像生成, gpt-image-2, Social Media Post, 社交媒体",
    "promptPreview": "Inside an elevator, the metal walls have a slight cold reflection, and the ceiling lights are whitish but uneven. The space is enclosed and quiet. A {argument n",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "social_poster",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-442",
    "locale": "zh",
    "title": "图像模板 - 社交媒体海报",
    "category": "图像模板 - 社交媒体海报",
    "description": "anime, fantasy, 图像生成, gpt-image-2, Social Media Post, 社交媒体",
    "promptPreview": "A cute pastel anime illustration of a young elf girl streamer or office worker sitting at a desk and typing on a mechanical keyboard in a cozy bedroom workspace",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "social_poster",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-443",
    "locale": "zh",
    "title": "图像模板 - 社交媒体海报",
    "category": "图像模板 - 社交媒体海报",
    "description": "图像生成, gpt-image-2, Social Media Post, 时尚, 社交媒体, 编辑风",
    "promptPreview": "A woman with {argument name=\"hair color\" default=\"long red hair\"} crouching in a minimalist studio setting with a {argument name=\"background color\" default=\"sof",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "social_poster",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-445",
    "locale": "zh",
    "title": "图像模板 - 社交媒体海报",
    "category": "图像模板 - 社交媒体海报",
    "description": "cinematic, typography, 图像生成, gpt-image-2, Social Media Post, 海报, 社交媒体",
    "promptPreview": "Create a dramatic football transfer announcement poster in a vertical social-media format, centered on a photorealistic adult male soccer player wearing a moder",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "social_poster",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-446",
    "locale": "zh",
    "title": "图像模板 - 社交媒体海报",
    "category": "图像模板 - 社交媒体海报",
    "description": "storyboard, dance, portrait, cinematic, sequence, fashion, 图像生成, gpt-image-2, Social Media Post, 社交媒体",
    "promptPreview": "# Sensational Girl Dance — 8-Shot Storyboard for GPT-Image-2 For each shot, prepend the GLOBAL STYLE TOKENS and append the NEGATIVE PROMPT. Shots are choreograp",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "social_poster",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-447",
    "locale": "zh",
    "title": "图像模板 - 社交媒体海报",
    "category": "图像模板 - 社交媒体海报",
    "description": "anime, typography, food, 图像生成, gpt-image-2, Social Media Post, 社交媒体",
    "promptPreview": "type: retro Japanese lifestyle magazine cover poster | theme: Showa Day feature celebrating nostalgic Japanese retro culture | style: clean editorial layout mix",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "social_poster",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-448",
    "locale": "zh",
    "title": "图像模板 - 社交媒体海报",
    "category": "图像模板 - 社交媒体海报",
    "description": "fantasy, 图像生成, gpt-image-2, Social Media Post, 时尚, 社交媒体",
    "promptPreview": "Based on this character info card for a {argument name=\"subject\" default=\"girl\"}, generate a 7-day outfit recommendation guide suitable for her appearance, heig",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "social_poster",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-449",
    "locale": "zh",
    "title": "图像模板 - 社交媒体海报",
    "category": "图像模板 - 社交媒体海报",
    "description": "cinematic, fantasy, 图像生成, gpt-image-2, Social Media Post, 旅行, 社交媒体, 拼图",
    "promptPreview": "A 12-frame collage of candid, emotional snapshots of a young {argument name=\"ethnicity\" default=\"Chinese\"} woman traveling alone in {argument name=\"location\" de",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "social_poster",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-450",
    "locale": "zh",
    "title": "图像模板 - 社交媒体海报",
    "category": "图像模板 - 社交媒体海报",
    "description": "action, 图像生成, gpt-image-2, Social Media Post, 社交媒体, 复古",
    "promptPreview": "A hand-lettered sketch of the phrase “{argument name=\"phrase\" default=\"Good Morning\"} ” on warm-white marker paper, drawn with a black brush marker. Soft graphi",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "social_poster",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-451",
    "locale": "zh",
    "title": "图像模板 - 社交媒体海报",
    "category": "图像模板 - 社交媒体海报",
    "description": "fantasy, 3d-render, product, 图像生成, gpt-image-2, Social Media Post, 海报, VR",
    "promptPreview": "type: exploded view product diagram poster | subject: VR headset | style: clean high-tech 3D render, studio lighting, glowing accents | background: {argument na",
    "sourceId": "davidwuw",
    "sourceLanguage": "ZH",
    "sourceLabel": "open-design",
    "tags": [
      "social_poster",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-75",
    "locale": "zh",
    "title": "名画人物麦片广告",
    "category": "广告设计",
    "description": "Upload a portrait image as reference.",
    "promptPreview": "\"Master Oats\": Based on the character features in the photo I uploaded, generate a suitable oatmeal pairing for them (such as vegetables, fruits, yogurt, whole",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "advertising",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-1",
    "locale": "zh",
    "title": "实物+手绘涂鸦创意广告",
    "category": "广告设计",
    "description": "Replace [real object], [doodle concept and interaction], [ad copy], and [brand logo] with specific content.",
    "promptPreview": "A minimalist and creative advertisement set on a pure white background. A real [real object] combined with hand-drawn black ink doodles, with loose and playful",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "advertising",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-64",
    "locale": "zh",
    "title": "柔和3D广告",
    "category": "广告设计",
    "description": "Replace [brand product] with a specific product description.",
    "promptPreview": "A pastel 3D cartoon-style [brand product] sculpture, made of smooth clay-like textures and vibrant pastel colors, placed in a minimalist isometric scene that co",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "advertising",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-21",
    "locale": "zh",
    "title": "代码风格名片",
    "category": "文字渲染",
    "description": "Replace the name, title, email, and link data in the JSON code.",
    "promptPreview": "Close-up shot: a hand holding a business card designed to look like a JSON file in VS Code. The code on the card is displayed with real JSON syntax highlighting",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "text_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-84",
    "locale": "zh",
    "title": "复古CRT电脑开机画面",
    "category": "文字渲染",
    "description": "Replace [shape or logo] with a specific description.",
    "promptPreview": "Retro CRT computer boot screen, finally displaying ASCII art in the shape of [shape or logo].",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "text_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-14",
    "locale": "zh",
    "title": "字母与词义融合",
    "category": "文字渲染",
    "description": "Replace { beautify } with the word you want to blend.",
    "promptPreview": "Incorporate the meaning of the word into its letters, cleverly blending graphics and letters together. Word: { beautify } Add a brief description of the word be",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "text_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-63",
    "locale": "zh",
    "title": "手绘信息图卡片",
    "category": "文字渲染",
    "description": "Add your own text content at the end of the prompt.",
    "promptPreview": "Create a hand-drawn style infographic card in 9:16 portrait ratio. The card has a clear theme, with a paper-textured beige or off-white background, and the over",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "text_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-67",
    "locale": "zh",
    "title": "手绘信息图卡片(认知)",
    "category": "文字渲染",
    "description": "Replace the text content at the end.",
    "promptPreview": "Create a hand-drawn style infographic card in 9:16 portrait ratio. The card has a clear theme, with a paper-textured beige or off-white background, and the over",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "text_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-50",
    "locale": "zh",
    "title": "护照入境印章",
    "category": "文字渲染",
    "description": "Replace city, country, landmark, and date in brackets.",
    "promptPreview": "Create a realistic passport page stamped with an entry stamp from [Beijing, China]. The stamp should have \"Welcome to Beijing\" written in bold English, designed",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "text_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-40",
    "locale": "zh",
    "title": "虚构推文截图(爱因斯坦)",
    "category": "文字渲染",
    "description": "Fictional Tweet Screenshot (Einstein)",
    "promptPreview": "A hyper-realistic tweet posted by Einstein right after completing the theory of relativity. It includes a selfie where the chalkboard and scrawled formulas are",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "text_render",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-90",
    "locale": "zh",
    "title": "PS2游戏封面",
    "category": "海报设计",
    "description": "PS2 Game Cover (GTA x Shrek)",
    "promptPreview": "Can you create an image of a PS2 game cover? The title is \"Grand Theft Auto: Far Far Away\". It's a GTA-style game set in the Shrek universe.",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "poster",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-98",
    "locale": "zh",
    "title": "复古宣传海报",
    "category": "海报设计",
    "description": "Retro Propaganda Poster",
    "promptPreview": "Vintage propaganda poster style, with prominent Chinese text, and a red-yellow radial pattern background. In the center of the image is a beautiful young woman,",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "poster",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-79",
    "locale": "zh",
    "title": "小红书封面",
    "category": "海报设计",
    "description": "Replace the copy content for different covers.",
    "promptPreview": "Draw an image: Create a Xiaohongshu (Little Red Book) cover. Requirements: Have enough appeal to attract users to click; Use eye-catching fonts, choose distinct",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "poster",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-42",
    "locale": "zh",
    "title": "彩色矢量艺术海报",
    "category": "海报设计",
    "description": "Replace city and country names to generate different city posters.",
    "promptPreview": "The location is \"London, United Kingdom,\" generate a colorful vector art poster featuring iconic landmarks, cultural symbols, and local elements arranged in a v",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "poster",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-52",
    "locale": "zh",
    "title": "时尚杂志封面风格",
    "category": "海报设计",
    "description": "Fashion Magazine Cover Style",
    "promptPreview": "A beautiful woman wearing a pink cheongsam (qipao), with an exquisite floral headpiece, her hair adorned with colorful flowers, and an elegant white lace collar",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "poster",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-92",
    "locale": "zh",
    "title": "极简未来主义海报",
    "category": "海报设计",
    "description": "Replace [3D Coca-Cola classic soda bottle] with other objects for different themes.",
    "promptPreview": "A vertical (3:4) 4K resolution minimalist futurist exhibition poster with an ultra-light cool gray #f4f4f4 background. At the center of the poster is a fluid 3D",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "poster",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-86",
    "locale": "zh",
    "title": "讽刺海报生成",
    "category": "海报设计",
    "description": "Replace with your own satirical theme.",
    "promptPreview": "Generate a satirical poster for me: GPT 4o is dominating so hard, everyone should quit image AI and go deliver food instead",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "poster",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-473",
    "locale": "zh",
    "title": "视频模板 - 动漫改编",
    "category": "视频模板 - 动漫改编",
    "description": "anime, fantasy, action, 视频生成, seedance-2.0",
    "promptPreview": "Live-Action Anime Adaptation · Breathing Technique Decisive Battle (15 seconds · Super Burning Special Effects Version) 【Core Focus】: Water Breathing (Blue Wate",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "anime_adaptation",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-452",
    "locale": "zh",
    "title": "视频模板 - 动画",
    "category": "视频模板 - 动画",
    "description": "视频生成, seedance-2.0, General",
    "promptPreview": "Scene: A boy in a room seriously assembling Lego blocks. The visual style is 3D animation with vibrant colors, smooth lines, full of childlike fun and vitality.",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "animation",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-456",
    "locale": "zh",
    "title": "视频模板 - 动画",
    "category": "视频模板 - 动画",
    "description": "视频生成, seedance-2.0, General",
    "promptPreview": "apply the walking animation of @anim excatly as it is to @char7 . the camera tracks the character exactly in place, camera angle does not change",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "animation",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-472",
    "locale": "zh",
    "title": "视频模板 - 动画",
    "category": "视频模板 - 动画",
    "description": "视频生成, seedance-2.0, General",
    "promptPreview": "create a walking animation for this hunched over character. the character stays in place",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "animation",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-477",
    "locale": "zh",
    "title": "视频模板 - 动画",
    "category": "视频模板 - 动画",
    "description": "cinematic, typography, 视频生成, seedance-2.0, General",
    "promptPreview": "Hyper-detailed nightclub event flyer, a Black woman with sleek long braids and bold red lipstick photographed in a confident chin-up power pose looking directly",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "animation",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-488",
    "locale": "zh",
    "title": "视频模板 - 动画",
    "category": "视频模板 - 动画",
    "description": "视频生成, seedance-2.0, General, 复古",
    "promptPreview": "Classic vintage Disney animation style. Scene 1: On a pirate ship sailing on the sea, a fat and sinister crocodile pirate stands at the end of a plank. Three bi",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "animation",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-453",
    "locale": "zh",
    "title": "视频模板 - 广告/品牌",
    "category": "视频模板 - 广告/品牌",
    "description": "cinematic, fantasy, product, 视频生成, seedance-2.0, Advertising",
    "promptPreview": "Create a 15-second ultra-realistic cinematic transformation video using the exact same man from the uploaded reference image. Maintain perfect face consistency,",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "advertising",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-475",
    "locale": "zh",
    "title": "视频模板 - 广告/品牌",
    "category": "视频模板 - 广告/品牌",
    "description": "cinematic, 视频生成, seedance-2.0, Advertising, 分镜",
    "promptPreview": "0s – 4s (Arrival at the Academy) A massive gothic magical academy appears above floating cliffs, surrounded by storm clouds and glowing runes. The girl walks th",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "advertising",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-484",
    "locale": "zh",
    "title": "视频模板 - 广告/品牌",
    "category": "视频模板 - 广告/品牌",
    "description": "cinematic, fantasy, action, 视频生成, seedance-2.0, Advertising, 舞蹈",
    "promptPreview": "Use the first reference image as the exact choreography and motion-process guide. Use the second reference image as the identity reference for the adult woman d",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "advertising",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-454",
    "locale": "zh",
    "title": "视频模板 - 武侠/历史",
    "category": "视频模板 - 武侠/历史",
    "description": "fantasy, action, 视频生成, seedance-2.0, General, 龙",
    "promptPreview": "Shot 1 (00:00–00:02) – WS, Rainy Night, Forward Tracking. A narrow, ancient village alley drenched in relentless rain. Water streams down slanted rooftops and f",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "wuxia_history",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-455",
    "locale": "zh",
    "title": "视频模板 - 武侠/历史",
    "category": "视频模板 - 武侠/历史",
    "description": "cinematic, nature, 视频生成, seedance-2.0, General",
    "promptPreview": "extremely fast-paced cinematic FPV flying through the ancient Indian Dandaka kingdom, dense mystical forests, towering sal and teak trees, tribal settlements, a",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "wuxia_history",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-469",
    "locale": "zh",
    "title": "视频模板 - 游戏/科幻",
    "category": "视频模板 - 游戏/科幻",
    "description": "cinematic, cyberpunk, 3d-render, 视频生成, seedance-2.0, General, 游戏",
    "promptPreview": "INK INDUSTRIES : GAME TRAILER CHARACTERYoung athletic male, dark curly hair, shirtless, full chest and back tattoos, gold hoop earring, cigarette in mouth, blac",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "game_scifi",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-471",
    "locale": "zh",
    "title": "视频模板 - 特效/奇幻",
    "category": "视频模板 - 特效/奇幻",
    "description": "cinematic, fantasy, 3d-render, 视频生成, seedance-2.0, VFX / Fantasy",
    "promptPreview": "[Style] Hollywood Haute Couture Fantasy blockbuster, 8K ultra-clear, Photorealistic, High-fashion Editorial Style, Unreal Engine 5 fluid rendering, visual illus",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "vfx_fantasy",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-482",
    "locale": "zh",
    "title": "视频模板 - 特效/奇幻",
    "category": "视频模板 - 特效/奇幻",
    "description": "cinematic, 视频生成, seedance-2.0, VFX / Fantasy",
    "promptPreview": "0s – 4s (Opening Mystery) A hidden magical kingdom under heavy rain at night. The girl stands beside a glowing water mirror in an ancient palace courtyard. Blue",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "vfx_fantasy",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-459",
    "locale": "zh",
    "title": "视频模板 - 电影叙事",
    "category": "视频模板 - 电影叙事",
    "description": "cinematic, fantasy, cinematic-romance, 视频生成, seedance-2.0, 生日",
    "promptPreview": "0s–4s Close-up of a young girl waking up in a softly lit bedroom, warm golden sunlight through curtains, she gently smiles while checking her phone filled with",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "cinematic",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-460",
    "locale": "zh",
    "title": "视频模板 - 电影叙事",
    "category": "视频模板 - 电影叙事",
    "description": "cinematic, fantasy, action, 视频生成, seedance-2.0, 龙",
    "promptPreview": "STYLE Handheld + aerial camera blend Soft motion blur (only during fast transitions) Teal–orange cinematic grade Cool tones during dragon moments, warm tones at",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "cinematic",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-461",
    "locale": "zh",
    "title": "视频模板 - 电影叙事",
    "category": "视频模板 - 电影叙事",
    "description": "cinematic, action, 视频生成, seedance-2.0, 舞蹈",
    "promptPreview": "1 0-3s Extreme close-up of the face, exquisite and three-dimensional features, cold and elegant eyes locked on the lens, sword dance opening pose: hands quickly",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "cinematic",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-462",
    "locale": "zh",
    "title": "视频模板 - 电影叙事",
    "category": "视频模板 - 电影叙事",
    "description": "portrait, cinematic, action, 视频生成, seedance-2.0",
    "promptPreview": "A realistic human face with highly detailed skin texture, pores, and micro-musculature. Scene: Tight portrait close-up against a dark, void-like background. Sty",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "cinematic",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-463",
    "locale": "zh",
    "title": "视频模板 - 电影叙事",
    "category": "视频模板 - 电影叙事",
    "description": "cinematic, 视频生成, seedance-2.0",
    "promptPreview": "A marine biologist in a sleek wetsuit swims through the vibrant coral reefs of the Great Barrier Reef. At the 3-second mark, he dives deeper to approach an anci",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "cinematic",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-464",
    "locale": "zh",
    "title": "视频模板 - 电影叙事",
    "category": "视频模板 - 电影叙事",
    "description": "cinematic, cyberpunk, fantasy, 视频生成, seedance-2.0, 音乐",
    "promptPreview": "**Cinematic Truth Source & Setup** Professional music podcast video production, shot on Sony FX6 cinema camera in 4K DCI, anamorphic lenses with natural breathi",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "cinematic",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-465",
    "locale": "zh",
    "title": "视频模板 - 电影叙事",
    "category": "视频模板 - 电影叙事",
    "description": "cinematic, fantasy, action, 视频生成, seedance-2.0, 导航",
    "promptPreview": "Create a 5-second cinematic route-guide clip for a walking navigation video. Continuity: This is scene {N} of 5 in a route from North Avenue MARTA Station to Co",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "cinematic",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-466",
    "locale": "zh",
    "title": "视频模板 - 电影叙事",
    "category": "视频模板 - 电影叙事",
    "description": "cinematic, cyberpunk, action, 视频生成, seedance-2.0, 舞蹈, 赛车",
    "promptPreview": "cinematic street racing sequence at night, a focused driver inside a high-performance car grips the steering wheel, intense eye focus, city lights reflecting on",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "cinematic",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-467",
    "locale": "zh",
    "title": "视频模板 - 电影叙事",
    "category": "视频模板 - 电影叙事",
    "description": "cinematic, cyberpunk, 视频生成, seedance-2.0",
    "promptPreview": "Draev stands in the center of the neon-lit alley, surrounded by multiple vampires positioned on rooftops and street level. The vampires attack simultaneously. D",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "cinematic",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-468",
    "locale": "zh",
    "title": "视频模板 - 电影叙事",
    "category": "视频模板 - 电影叙事",
    "description": "cinematic, 视频生成, seedance-2.0",
    "promptPreview": "SHOT 1: Cinematic wide angle format — rocket launching into night sky, city lights below, clouds parting, stars above. Dark, dramatic, photorealistic. SHOT 2: M",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "cinematic",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-474",
    "locale": "zh",
    "title": "视频模板 - 电影叙事",
    "category": "视频模板 - 电影叙事",
    "description": "cinematic, action, nature, 视频生成, seedance-2.0",
    "promptPreview": "Global Intent: Quiet Luxury with an aggressive edge. A stylish man with Dobermans and a classic dark blue vintage supercar journeys through misty mountains to a",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "cinematic",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-476",
    "locale": "zh",
    "title": "视频模板 - 电影叙事",
    "category": "视频模板 - 电影叙事",
    "description": "cinematic, action, food, 视频生成, seedance-2.0, 乡村, 治愈",
    "promptPreview": "[Style] Modern Rural Aesthetics, Cinematic Commercial quality, shot with Sony A7S3/cinema camera, 4K/8K ultra-clear, Extreme Macro, natural transparent lighting",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "cinematic",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-478",
    "locale": "zh",
    "title": "视频模板 - 电影叙事",
    "category": "视频模板 - 电影叙事",
    "description": "cinematic, product, 视频生成, seedance-2.0, 武侠",
    "promptPreview": "core_style: 80s-90s Shaw Brothers film style, early Hong Kong Wuxia drama aesthetics, nostalgic Chinese Wuxia movies, vintage TV quality, warm tones with high s",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "cinematic",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-485",
    "locale": "zh",
    "title": "视频模板 - 电影叙事",
    "category": "视频模板 - 电影叙事",
    "description": "game-cinematic, arpg, three-kingdoms, ancient-china, combat, cavalry, guanyu, key-visual, hud-safe, companion-to-image",
    "promptPreview": "A ~10 second in-engine cinematic ARPG action sequence, photoreal, Unreal-Engine-5-grade render quality, desaturated filmic color grading, shallow depth of field",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "cinematic",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-486",
    "locale": "zh",
    "title": "视频模板 - 电影叙事",
    "category": "视频模板 - 电影叙事",
    "description": "game-cinematic, arpg, three-kingdoms, ancient-china, combat, archery, lyubu, key-visual, hud-safe, companion-to-image",
    "promptPreview": "A ~10 second in-engine cinematic ARPG action sequence, photoreal, Unreal-Engine-5-grade render quality, desaturated filmic color grading, shallow depth of field",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "cinematic",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-487",
    "locale": "zh",
    "title": "视频模板 - 电影叙事",
    "category": "视频模板 - 电影叙事",
    "description": "game-cinematic, arpg, three-kingdoms, ancient-china, combat, cavalry, zhaoyun, escort, key-visual, hud-safe",
    "promptPreview": "A ~12 second in-engine cinematic ARPG action sequence, photoreal, Unreal-Engine-5-grade render quality, desaturated filmic color grading, shallow depth of field",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "cinematic",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-458",
    "locale": "zh",
    "title": "视频模板 - 短视频/影像",
    "category": "视频模板 - 短视频/影像",
    "description": "cinematic, fantasy, 3d-render, 视频生成, seedance-2.0, Motion Graphics",
    "promptPreview": "Based on the three characters in the reference images. High definition, Unreal Engine rendering, cinematic quality, candy-colored palette, Japanese-style aesthe",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "short_video",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-470",
    "locale": "zh",
    "title": "视频模板 - 短视频/影像",
    "category": "视频模板 - 短视频/影像",
    "description": "action, 视频生成, seedance-2.0, General",
    "promptPreview": "STORY FORMAT: 15s / 150 BPM / MULTI-CUT / American dark comedy with exaggerated imperial satire / slapstick timing and punchline ending TONE: tense accusation →",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "short_video",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-481",
    "locale": "zh",
    "title": "视频模板 - 短视频/影像",
    "category": "视频模板 - 短视频/影像",
    "description": "fantasy, 3d-render, 视频生成, seedance-2.0, General, 武术",
    "promptPreview": "[STYLE] Monochrome grayscale illustration, 3D-rendered character, clean instructional reference sheet, white background, comic-style cell grid layout, technical",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "short_video",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-483",
    "locale": "zh",
    "title": "视频模板 - 短视频/影像",
    "category": "视频模板 - 短视频/影像",
    "description": "food, 视频生成, seedance-2.0, General",
    "promptPreview": "A realistic shot of an old man in a cozy kitchen being jumpscared when his toaster launches the bread five feet into the air like a rocket. Handheld \"home video",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "short_video",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-490",
    "locale": "zh",
    "title": "视频模板 - 短视频/影像",
    "category": "视频模板 - 短视频/影像",
    "description": "action, 视频生成, seedance-2.0, General",
    "promptPreview": "Ultra-realistic desert horizon. A gigantic industrial factory moving on mechanical legs crosses the wasteland like a living city. Female rebel riding a fast bik",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "short_video",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-489",
    "locale": "zh",
    "title": "视频模板 - 社交/舞蹈",
    "category": "视频模板 - 社交/舞蹈",
    "description": "fantasy, 3d-render, 视频生成, seedance-2.0, Social / Meme, 舞蹈, K-pop",
    "promptPreview": "viral kpop dance. Monochrome grayscale illustration, 3D-rendered character, clean instructional reference sheet, white background, comic-style cell grid layout,",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "social_dance",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-457",
    "locale": "zh",
    "title": "视频模板 - 舞蹈/动作",
    "category": "视频模板 - 舞蹈/动作",
    "description": "fantasy, 视频生成, seedance-2.0, General, 舞蹈",
    "promptPreview": "Have the character from Image 1 perform the dance based on the breakdown in Image 3. During the performance, include a beat-synced transformation into the chara",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "dance_action",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-480",
    "locale": "zh",
    "title": "视频模板 - 舞蹈/动作",
    "category": "视频模板 - 舞蹈/动作",
    "description": "cyberpunk, action, 视频生成, seedance-2.0, General, 舞蹈",
    "promptPreview": "16:9 horizontal screen, street rap MV style, neon purple and blue cool tones, explosive cool and fierce atmosphere. 0-3 seconds: Medium shot push-in, city stree",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "open-design",
    "tags": [
      "dance_action",
      "open-design"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-493",
    "locale": "zh",
    "title": "刘裕 - 山海传说角色定帧",
    "category": "角色肖像",
    "description": "山海传说（Legend of the Frontier）角色定帧图，刘裕，老板原创",
    "promptPreview": "A powerful portrait of Liu Yu, the future emperor. He sits on a heavy iron throne. Behind him is a giant, weathered bronze shield with dragon engravings, toweri",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "original",
    "tags": [
      "character-portrait",
      "original"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-494",
    "locale": "zh",
    "title": "拓跋珪 - 山海传说角色定帧",
    "category": "角色肖像",
    "description": "山海传说（Legend of the Frontier）角色定帧图，拓跋珪，老板原创",
    "promptPreview": "A dominant portrait of Tuoba Gui, the Northern Wei founder. He sits on a throne made of dark basalt. Behind him is a colossal, intimidating stone sculpture of a",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "original",
    "tags": [
      "character-portrait",
      "original"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-491",
    "locale": "zh",
    "title": "燕飞 - 山海传说角色定帧",
    "category": "角色肖像",
    "description": "山海传说（Legend of the Frontier）角色定帧图，燕飞，老板原创",
    "promptPreview": "A majestic cinematic portrait of Yan Fei from \"Legend of the Frontier\". He is sitting on a minimalist stone throne. Behind him is a colossal, ethereal crystalli",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "original",
    "tags": [
      "character-portrait",
      "original"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-492",
    "locale": "zh",
    "title": "纪千千 - 山海传说角色定帧",
    "category": "角色肖像",
    "description": "山海传说（Legend of the Frontier）角色定帧图，纪千千，老板原创",
    "promptPreview": "A breathtakingly beautiful cinematic portrait of Ji Qianqian, the most legendary beauty of the Southern Dynasties. She is seated gracefully on a delicate, carve",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "original",
    "tags": [
      "character-portrait",
      "original"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-71",
    "locale": "zh",
    "title": "35mm胶片风飞行岛",
    "category": "风景/场景",
    "description": "Replace city name for different scenes.",
    "promptPreview": "35mm film-style photo: Moscow floating on a flying island in the sky.",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "scene",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-7",
    "locale": "zh",
    "title": "三只动物地标自拍",
    "category": "风景/场景",
    "description": "Replace [animal type] and [landmark] with specific descriptions.",
    "promptPreview": "A close-up selfie photo of three [animal type] in front of the iconic [landmark], each with a different expression, shot during golden hour with cinematic light",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "scene",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-94",
    "locale": "zh",
    "title": "个性化房间设计",
    "category": "风景/场景",
    "description": "Customize room elements as needed.",
    "promptPreview": "Generate my room design (bed, bookshelf, sofa, plants, computer desk and computer), with paintings hanging on the wall, and a city night view outside the window",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "scene",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-22",
    "locale": "zh",
    "title": "乐高城市景观",
    "category": "风景/场景",
    "description": "Can reference this prompt to generate other city landscapes.",
    "promptPreview": "Create a highly detailed and vibrantly colored LEGO version of the Shanghai Bund scene. The foreground features the classic historic Bund buildings, exquisitely",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "scene",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-43",
    "locale": "zh",
    "title": "云朵艺术",
    "category": "风景/场景",
    "description": "Replace [SUBJECT/OBJECT] and [LOCATION]. E.g., Chinese dragon above the Great Wall.",
    "promptPreview": "Generate a photo: capturing a daytime scene where scattered clouds in the sky form the shape of [SUBJECT/OBJECT], positioned above [LOCATION].",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "scene",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-82",
    "locale": "zh",
    "title": "国家微缩玩具箱",
    "category": "风景/场景",
    "description": "Replace [Country Name] with the specific country name.",
    "promptPreview": "A hyper-realistic overhead photograph showing a 3D-printed diorama inside a beige cardboard box, with the lid held open by two human hands. The interior of the",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "scene",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-19",
    "locale": "zh",
    "title": "城市主题天气预报",
    "category": "风景/场景",
    "description": "City, weather, temperature, and building names can be customized.",
    "promptPreview": "From a clear 45-degree overhead angle, display an isometric miniature model scene featuring iconic city buildings such as [Shanghai Oriental Pearl Tower, the Bu",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "scene",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-60",
    "locale": "zh",
    "title": "微缩场景模型",
    "category": "风景/场景",
    "description": "Replace scene in brackets with other Chinese stories.",
    "promptPreview": "Miniature diorama scene presentation, using tilt-shift photography techniques, depicting a chibi scene of [Sun Wukong's Three Battles with the White Bone Demon]",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "scene",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-30",
    "locale": "zh",
    "title": "赛博朋克微缩移轴景观",
    "category": "风景/场景",
    "description": "Replace 【Cyberpunk】 with other styles like futuristic city, steampunk, medieval village, etc.",
    "promptPreview": "An ultra-detailed miniature 【Cyberpunk】 landscape viewed from above, using a tilt-shift lens effect. The scene is filled with toy-like elements, all rendered in",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "scene",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-38",
    "locale": "zh",
    "title": "Emoji奶油冰淇淋",
    "category": "食物摄影",
    "description": "Replace 【🍓】 with other Emojis.",
    "promptPreview": "Generate an image: Transform 【🍓】 into a cream ice cream bar, with cream flowing in curves on top of the ice cream looking delicious and appetizing, floating at",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "food",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "davidwuw-36",
    "locale": "zh",
    "title": "梦幻海底场景冰棍",
    "category": "食物摄影",
    "description": "Fantasy Underwater Scene Popsicle",
    "promptPreview": "A tilted first-person perspective shot of a hand holding a surreal popsicle. The popsicle has a transparent blue shell, inside which an underwater scene is disp",
    "sourceId": "davidwuw",
    "sourceLanguage": "EN",
    "sourceLabel": "moosl/awsome-gpt-image-2-prompts",
    "tags": [
      "food",
      "moosl/awsome-gpt-image-2-prompts"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-402",
    "locale": "zh",
    "title": "3D 小红书个人资料卡",
    "category": "UI 与界面",
    "description": "一只手中握着一张3D小红书个人资料卡，卡片中间方形镂空，一个女孩随意地坐在卡片镂空的边缘，温暖的米色和柔和的粉彩美学背景，逼真的深度和阴影，电影般的柔和光线，闪亮光滑的纹理，小红书风格的UI，漂浮的互动图标（点赞、评论、分享）带有发光的霓虹效果，闪光和光晕，背景中温馨的美学布置包",
    "promptPreview": "一只手中握着一张3D小红书个人资料卡，卡片中间方形镂空，一个女孩随意地坐在卡片镂空的边缘，温暖的米色和柔和的粉彩美学背景，逼真的深度和阴影，电影般的柔和光线，闪亮光滑的纹理，小红书风格的UI，漂浮的互动图标（点赞、评论、分享）带有发光的霓虹效果，闪光和光晕，背景中温馨的美学布置包括书籍、花瓶里的花和一台复古相机，梦幻氛",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrGafish",
    "tags": [
      "UI & Interfaces",
      "UI",
      "3D",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-387",
    "locale": "en",
    "title": "Netflix 首页主视觉 UI",
    "category": "UI 与界面",
    "description": "Create a Netflix homepage UI featuring a main hero film with its title and still generated from the uploaded reference.",
    "promptPreview": "Create a Netflix homepage UI featuring a main hero film with its title and still generated from the uploaded reference.",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@aimikoda",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-287",
    "locale": "zh",
    "title": "不知火舞的小红书主页",
    "category": "UI 与界面",
    "description": "[中文] 生成不知火舞的小红书主页截图 [English] Generate a screenshot of Mai Shiranui's Xiaohongshu homepage",
    "promptPreview": "[中文] 生成不知火舞的小红书主页截图 [English] Generate a screenshot of Mai Shiranui's Xiaohongshu homepage",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@rionaifantasy",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-336",
    "locale": "zh",
    "title": "个人网页视觉设计",
    "category": "UI 与界面",
    "description": "原文未公开，案例目标是生成一张高完成度的个人主页视觉设计图。",
    "promptPreview": "原文未公开，案例目标是生成一张高完成度的个人主页视觉设计图。",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "苍何原创实测（公众号文章《我逆向了 329 条 GPT-Image2 提示词模板，全部开源！》）",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-96",
    "locale": "en",
    "title": "主题海报版式设计",
    "category": "UI 与界面",
    "description": "{ \"type\": \"VTuber stream thumbnail\", \"character\": { \"hair\": \"long blonde twin tails with pink gradient ends\", \"eyes\": \"large pink anime eyes",
    "promptPreview": "{ \"type\": \"VTuber stream thumbnail\", \"character\": { \"hair\": \"long blonde twin tails with pink gradient ends\", \"eyes\": \"large pink anime eyes\", \"expression\": \"ch",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@sayaka_aiart",
    "tags": [
      "UI & Interfaces",
      "Poster",
      "Character",
      "3D",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-90",
    "locale": "en",
    "title": "信息图可视化设计",
    "category": "UI 与界面",
    "description": "GPT-Image-2 prompt: please automatically generate a top-tier concept poster / infographic-style movie poster centered around {argument name=",
    "promptPreview": "GPT-Image-2 prompt: please automatically generate a top-tier concept poster / infographic-style movie poster centered around {argument name=\"theme\" default=\"ran",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@A9Quant",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-82",
    "locale": "en",
    "title": "信息图可视化设计",
    "category": "UI 与界面",
    "description": "{ \"type\": \"scientific optical setup diagram\", \"main_setup\": { \"base\": \"optical breadboard table with grid of mounting holes\", \"beam\": \"red l",
    "promptPreview": "{ \"type\": \"scientific optical setup diagram\", \"main_setup\": { \"base\": \"optical breadboard table with grid of mounting holes\", \"beam\": \"red laser beam passing ho",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@HumanOS_v2",
    "tags": [
      "UI & Interfaces",
      "Infographic",
      "Realistic",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-42",
    "locale": "en",
    "title": "写实摄影风格图",
    "category": "UI 与界面",
    "description": "Express {argument name=\"subject\" default=\"a powerful AI builder\"} in a graffiti sketch style, presenting an overall visual effect of quick o",
    "promptPreview": "Express {argument name=\"subject\" default=\"a powerful AI builder\"} in a graffiti sketch style, presenting an overall visual effect of quick outlines, free deform",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@blanplan",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Realistic",
      "Illustration",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-239",
    "locale": "zh",
    "title": "刘亦菲抖音直播畅聊中",
    "category": "UI 与界面",
    "description": "[中文] 9:16 的图片比例，生成一张抖音直播的截图，里面是 刘亦菲 在直播，刘亦菲 手里拿着牌子，牌子里写着 今晚直播，欢迎来参亦菲畅聊！ [English] 9:16 aspect ratio, generate a screenshot of a Douyin live ",
    "promptPreview": "[中文] 9:16 的图片比例，生成一张抖音直播的截图，里面是 刘亦菲 在直播，刘亦菲 手里拿着牌子，牌子里写着 今晚直播，欢迎来参亦菲畅聊！ [English] 9:16 aspect ratio, generate a screenshot of a Douyin live stream, inside is Li",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@alanblogsooo",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-177",
    "locale": "zh",
    "title": "吉利银河暗黑中控界面",
    "category": "UI 与界面",
    "description": "[中文] 帮我生成一个吉利银河m9的中控界面，尺寸为21:9，暗色系 [English] Help me generate a central control interface of Geely Galaxy M9, size 21:9, dark color scheme.",
    "promptPreview": "[中文] 帮我生成一个吉利银河m9的中控界面，尺寸为21:9，暗色系 [English] Help me generate a central control interface of Geely Galaxy M9, size 21:9, dark color scheme.",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@xin_pai88825",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-95",
    "locale": "en",
    "title": "品牌视觉识别图",
    "category": "UI 与界面",
    "description": "{ \"type\": \"anime-style livestream thumbnail\", \"character\": { \"hair\": \"{argument name=\\\"hair color\\\" default=\\\"short silver hair with cyan un",
    "promptPreview": "{ \"type\": \"anime-style livestream thumbnail\", \"character\": { \"hair\": \"{argument name=\\\"hair color\\\" default=\\\"short silver hair with cyan underlights\\\"}\", \"eyes",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@sayaka_aiart",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Poster",
      "Brand",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-227",
    "locale": "zh",
    "title": "哔哩哔哩户晨风直播截图",
    "category": "UI 与界面",
    "description": "[中文] 9:16 的图片，生成一张哔哩哔哩直播的截图，里面是 户晨风在直播，户晨风表情开心，手里拿着牌子，牌子里写着 “Austin总太性情了，大家给Austin总点点关注。” [English] A 9:16 image, generate a screenshot of a",
    "promptPreview": "[中文] 9:16 的图片，生成一张哔哩哔哩直播的截图，里面是 户晨风在直播，户晨风表情开心，手里拿着牌子，牌子里写着 “Austin总太性情了，大家给Austin总点点关注。” [English] A 9:16 image, generate a screenshot of a Bilibili live strea",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@austinit",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-400",
    "locale": "zh",
    "title": "多风格签名选择海报",
    "category": "UI 与界面",
    "description": "你是一个高端签名设计系统 + 风格人格视觉系统。 任务： 仅基于用户输入的「姓名」，生成一张「多风格签名选择海报（卡片式结构）」。 目标是把名字转译为具有笔势、气质与力量感的签名设计系统，让用户产生选择欲、认同感和分享欲。 输入信息： 姓名：[输入你的昵称] 禁止要求额外信息，必",
    "promptPreview": "你是一个高端签名设计系统 + 风格人格视觉系统。 任务： 仅基于用户输入的「姓名」，生成一张「多风格签名选择海报（卡片式结构）」。 目标是把名字转译为具有笔势、气质与力量感的签名设计系统，让用户产生选择欲、认同感和分享欲。 输入信息： 姓名：[输入你的昵称] 禁止要求额外信息，必须自动完成气质与风格推断。 隐藏执行逻辑",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "GitHub prompt",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Poster",
      "Brand",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-243",
    "locale": "zh",
    "title": "定制专属风格界面设计系统",
    "category": "UI 与界面",
    "description": "[中文] 用xx风格帮我生成一套UI设计系统，包含网页、移动端、卡片、控件、按钮 以及其它 [English] Generate a UI design system for me in xx style, including web pages, mobile, cards, ",
    "promptPreview": "[中文] 用xx风格帮我生成一套UI设计系统，包含网页、移动端、卡片、控件、按钮 以及其它 [English] Generate a UI design system for me in xx style, including web pages, mobile, cards, controls, buttons, a",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@stark_nico99",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-280",
    "locale": "zh",
    "title": "封面排版设计图",
    "category": "UI 与界面",
    "description": "[中文] 以涂鸦速写风表现【一个厉害的AI builder】，整体呈现快速勾勒、自由变形、即兴手绘与草稿式的视觉效果。线条随手、夸张、可粗细不一，略显凌乱但具有节奏和表现力，强调概括、夸张、趣味和随性，而不是严谨写实或精细刻画。 颜色采用粗糙、干刷感明显的块面表现，可保留不均匀的",
    "promptPreview": "[中文] 以涂鸦速写风表现【一个厉害的AI builder】，整体呈现快速勾勒、自由变形、即兴手绘与草稿式的视觉效果。线条随手、夸张、可粗细不一，略显凌乱但具有节奏和表现力，强调概括、夸张、趣味和随性，而不是严谨写实或精细刻画。 颜色采用粗糙、干刷感明显的块面表现，可保留不均匀的涂抹痕迹、刷痕、飞白与覆盖感，色彩根据【",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "OpenNana",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-403",
    "locale": "zh",
    "title": "小红书数字破屏 3D 女孩",
    "category": "UI 与界面",
    "description": "将一位气质绝佳的女孩放置在一个显示小红书帖文的3D透明玻璃手机画面中，并重新调整她的身体姿势，使她看起来像是正从屏幕中突破、进入现实世界。其中一只脚必须强烈地朝向观者延伸，采用戏剧化的 3D 透视效果，创造出强烈的深度感与沉浸感。整体姿势需要具有动态感、自然且符合人体结构，就像是",
    "promptPreview": "将一位气质绝佳的女孩放置在一个显示小红书帖文的3D透明玻璃手机画面中，并重新调整她的身体姿势，使她看起来像是正从屏幕中突破、进入现实世界。其中一只脚必须强烈地朝向观者延伸，采用戏剧化的 3D 透视效果，创造出强烈的深度感与沉浸感。整体姿势需要具有动态感、自然且符合人体结构，就像是在从屏幕中跨步而出的瞬间。 手机屏幕边缘",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrGafish",
    "tags": [
      "UI & Interfaces",
      "Character",
      "3D",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-161",
    "locale": "en",
    "title": "应用界面样机图",
    "category": "UI 与界面",
    "description": "{ \"type\": \"video game screenshot mockup\", \"perspective\": \"third-person over-the-shoulder\", \"character\": { \"description\": \"male protagonist s",
    "promptPreview": "{ \"type\": \"video game screenshot mockup\", \"perspective\": \"third-person over-the-shoulder\", \"character\": { \"description\": \"male protagonist seen from behind\", \"c",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@DanDaniDaniel01",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Character",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-156",
    "locale": "zh",
    "title": "应用界面样机图",
    "category": "UI 与界面",
    "description": "{ \"type\": \"mobile live-streaming e-commerce interface mockup\", \"subject\": { \"description\": \"young Asian woman, long dark hair, wearing light",
    "promptPreview": "{ \"type\": \"mobile live-streaming e-commerce interface mockup\", \"subject\": { \"description\": \"young Asian woman, long dark hair, wearing light-colored floral paja",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@linxiaobei888",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Product",
      "Character",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-107",
    "locale": "en",
    "title": "应用界面样机图",
    "category": "UI 与界面",
    "description": "{\"type\": \"YouTube desktop dark mode UI mockup\", \"header\": {\"logo\": \"YouTube\", \"search_bar\": \"Search\", \"icons_count\": 5, \"icons\": [\"search\", ",
    "promptPreview": "{\"type\": \"YouTube desktop dark mode UI mockup\", \"header\": {\"logo\": \"YouTube\", \"search_bar\": \"Search\", \"icons_count\": 5, \"icons\": [\"search\", \"mic\", \"create\", \"no",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@tehno_maniak",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Realistic",
      "Illustration",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-7",
    "locale": "zh",
    "title": "应用界面样机图",
    "category": "UI 与界面",
    "description": "生成一张竖版手机截图风格的图片，整体比例接近 9:16。画面中心偏上是一位真人 coser，扮演上传图片中的二次元角色。人物为写实风格，但五官略带动漫感，皮肤细腻，眼睛稍大，表情温柔地看向镜头，坐在室内的休闲场景中，例如咖啡厅或酒吧吧台前，背景有符合场景的道具。画面最上方加入手机",
    "promptPreview": "生成一张竖版手机截图风格的图片，整体比例接近 9:16。画面中心偏上是一位真人 coser，扮演上传图片中的二次元角色。人物为写实风格，但五官略带动漫感，皮肤细腻，眼睛稍大，表情温柔地看向镜头，坐在室内的休闲场景中，例如咖啡厅或酒吧吧台前，背景有符合场景的道具。画面最上方加入手机系统状态栏 UI，包括时间、电量、信号、",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "小红书号944846927",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Realistic",
      "Character",
      "Fashion",
      "Food",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-258",
    "locale": "zh",
    "title": "快手直播离婚预告手机截图",
    "category": "UI 与界面",
    "description": "[中文] 生成快手内容截图：主题：直播离婚预告，iPhone尺寸 [English] Generate Kuaishou content screenshot: Theme: Live divorce announcement, iPhone size",
    "promptPreview": "[中文] 生成快手内容截图：主题：直播离婚预告，iPhone尺寸 [English] Generate Kuaishou content screenshot: Theme: Live divorce announcement, iPhone size",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-440",
    "locale": "en",
    "title": "手机拍摄 FaceTime 工作屏幕",
    "category": "UI 与界面",
    "description": "Create a raw smartphone photo of a laptop screen, not a screenshot. Aspect ratio 3:4, high-angle downward POV from someone standing over a d",
    "promptPreview": "Create a raw smartphone photo of a laptop screen, not a screenshot. Aspect ratio 3:4, high-angle downward POV from someone standing over a desk at night. The la",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@kaanakz",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Realistic",
      "3D",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-257",
    "locale": "zh",
    "title": "抖音汉服美女直播带货截图",
    "category": "UI 与界面",
    "description": "[中文] 生成一个抖音直播的截图里面是一个穿着中国传统服饰的美女在直播卖货 [English] Generate a screenshot of a Douyin live stream, featuring a beautiful woman wearing tradition",
    "promptPreview": "[中文] 生成一个抖音直播的截图里面是一个穿着中国传统服饰的美女在直播卖货 [English] Generate a screenshot of a Douyin live stream, featuring a beautiful woman wearing traditional Chinese clothing ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@joshesye",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-308",
    "locale": "zh",
    "title": "抖音直播截图画面",
    "category": "UI 与界面",
    "description": "[中文] 9:16 的图片比例，生成一张抖音直播的截图，里面是 xxx 在直播，xxx 手里拿着牌子，牌子里写着 xxxx。 [English] 9:16 aspect ratio, generate a screenshot of a Douyin live stream, i",
    "promptPreview": "[中文] 9:16 的图片比例，生成一张抖音直播的截图，里面是 xxx 在直播，xxx 手里拿着牌子，牌子里写着 xxxx。 [English] 9:16 aspect ratio, generate a screenshot of a Douyin live stream, inside is xxx live st",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@_FORAB",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-256",
    "locale": "zh",
    "title": "抖音直播间的绝美女主播",
    "category": "UI 与界面",
    "description": "[中文] 生成一个抖音直播的截图 里面是一个美女在直播 [English] Generate a screenshot of a Douyin livestream, inside there is a beautiful woman livestreaming",
    "promptPreview": "[中文] 生成一个抖音直播的截图 里面是一个美女在直播 [English] Generate a screenshot of a Douyin livestream, inside there is a beautiful woman livestreaming",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@joshesye",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-288",
    "locale": "zh",
    "title": "抖音美女直播间界面设计",
    "category": "UI 与界面",
    "description": "[中文] 生成抖音直播间界面，内容是一个美女在直播 [English] Generate a TikTok live stream interface, the content is a beautiful woman live streaming",
    "promptPreview": "[中文] 生成抖音直播间界面，内容是一个美女在直播 [English] Generate a TikTok live stream interface, the content is a beautiful woman live streaming",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@msjiaozhu",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-269",
    "locale": "zh",
    "title": "拒绝盲目催婚的暖心视频号截图",
    "category": "UI 与界面",
    "description": "[中文] 生成视频号内容截图，主题：中老年不要盲目催婚，iPhone尺寸 [English] Generate a screenshot of WeChat Channels content, theme: middle-aged and elderly people shoul",
    "promptPreview": "[中文] 生成视频号内容截图，主题：中老年不要盲目催婚，iPhone尺寸 [English] Generate a screenshot of WeChat Channels content, theme: middle-aged and elderly people should not blindly urge m",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-93",
    "locale": "en",
    "title": "插画艺术风格创作",
    "category": "UI 与界面",
    "description": "{ \"type\": \"VTuber stream thumbnail\", \"style\": \"anime, highly detailed, cute, sparkly, overwhelmingly pink color palette\", \"character\": { \"de",
    "promptPreview": "{ \"type\": \"VTuber stream thumbnail\", \"style\": \"anime, highly detailed, cute, sparkly, overwhelmingly pink color palette\", \"character\": { \"description\": \"anime g",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@sayaka_aiart",
    "tags": [
      "UI & Interfaces",
      "Illustration",
      "Character",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-261",
    "locale": "zh",
    "title": "智能视频生成器暗黑界面设计",
    "category": "UI 与界面",
    "description": "[中文] 渲染一个专业的IOS APP首页UI图，该主题为AI Video Generator,英文界面。专业级设计，专业风格，暗黑色主题。 [English] Render a professional iOS APP homepage UI image, the theme ",
    "promptPreview": "[中文] 渲染一个专业的IOS APP首页UI图，该主题为AI Video Generator,英文界面。专业级设计，专业风格，暗黑色主题。 [English] Render a professional iOS APP homepage UI image, the theme is AI Video Generato",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@austinit",
    "tags": [
      "UI & Interfaces",
      "UI",
      "3D",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-188",
    "locale": "zh",
    "title": "暗黑极简头像网站视觉设计",
    "category": "UI 与界面",
    "description": "[中文] 用 ABCD（a black cover design) 的风格，为 图你太美 设计一个 vi 系统。图你太美是一个头像美图分享 网站。 [English] In the style of ABCD (a black cover design), design a VI",
    "promptPreview": "[中文] 用 ABCD（a black cover design) 的风格，为 图你太美 设计一个 vi 系统。图你太美是一个头像美图分享 网站。 [English] In the style of ABCD (a black cover design), design a VI system for Tu Ni Ta",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@xiaoxiaodong01",
    "tags": [
      "UI & Interfaces",
      "Poster",
      "Realistic",
      "Character",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-335",
    "locale": "zh",
    "title": "朋友圈截图生成",
    "category": "UI 与界面",
    "description": "原文未公开，重点展示 GPT-Image2 在高仿社交截图与中文排版场景中的能力。",
    "promptPreview": "原文未公开，重点展示 GPT-Image2 在高仿社交截图与中文排版场景中的能力。",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "苍何原创实测（公众号文章《我逆向了 329 条 GPT-Image2 提示词模板，全部开源！》）",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Social",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-282",
    "locale": "zh",
    "title": "温柔治愈系二次元手机截图",
    "category": "UI 与界面",
    "description": "[中文] 生成一张竖版手机截图风格的图片，整体比例接近 9:16。画面中心偏上是一位真人 coser，扮演（角色名称）的二次元角色。人物为写实风格，但五官略带动漫感，皮肤细腻，眼睛稍大，表情温柔地看向镜头，坐在室内的休闲场景中，例如咖啡厅或酒吧吧台前，背景有符合场景的道具。画面最",
    "promptPreview": "[中文] 生成一张竖版手机截图风格的图片，整体比例接近 9:16。画面中心偏上是一位真人 coser，扮演（角色名称）的二次元角色。人物为写实风格，但五官略带动漫感，皮肤细腻，眼睛稍大，表情温柔地看向镜头，坐在室内的休闲场景中，例如咖啡厅或酒吧吧台前，背景有符合场景的道具。画面最上方加入手机系统状态栏 UI，包括时间、",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@Zoulinshen",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Realistic",
      "Character",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-91",
    "locale": "en",
    "title": "游戏界面截图",
    "category": "UI 与界面",
    "description": "A highly detailed, realistic first-person video game screenshot of a next-generation voxel-based world. At the top center, a large, bold 3D ",
    "promptPreview": "A highly detailed, realistic first-person video game screenshot of a next-generation voxel-based world. At the top center, a large, bold 3D logo reads \"{argumen",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@wolfaidev",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Realistic",
      "Brand",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-164",
    "locale": "zh",
    "title": "特朗普太空直播间破千万",
    "category": "UI 与界面",
    "description": "[中文] 一张9:16竖屏的抖音直播截图，太空直播风格。特朗普穿着NASA风格的白色宇航服，头盔面罩半开，露出他标志性的金色头发和笑容。他漂浮在国际空间站的舱内进行直播，处于微重力失重状态，身体微微悬浮。他双手举着一块固定在宇航服上的金属铭牌，铭牌上用NASA风格的印刷体写着\"感",
    "promptPreview": "[中文] 一张9:16竖屏的抖音直播截图，太空直播风格。特朗普穿着NASA风格的白色宇航服，头盔面罩半开，露出他标志性的金色头发和笑容。他漂浮在国际空间站的舱内进行直播，处于微重力失重状态，身体微微悬浮。他双手举着一块固定在宇航服上的金属铭牌，铭牌上用NASA风格的印刷体写着\"感谢松果先森送的大火箭\"。身后圆形舷窗外可",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@songguoxiansen",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Brand",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-255",
    "locale": "zh",
    "title": "瑜伽裤女主播展示身材曲线",
    "category": "UI 与界面",
    "description": "[中文] 手机竖屏界面，短视频直播平台风格，一位年轻亚洲女主播在家中直播带货，主播穿着贴身瑜伽裤与简约上衣，身材曲线自然，正在侧身展示裤子的线条与弹性，动作自然不夸张； [English] Mobile vertical screen interface, short video",
    "promptPreview": "[中文] 手机竖屏界面，短视频直播平台风格，一位年轻亚洲女主播在家中直播带货，主播穿着贴身瑜伽裤与简约上衣，身材曲线自然，正在侧身展示裤子的线条与弹性，动作自然不夸张； [English] Mobile vertical screen interface, short video live streaming plat",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@joshesye",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-159",
    "locale": "zh",
    "title": "界面交互设计图",
    "category": "UI 与界面",
    "description": "{ \"type\": \"e-commerce livestream UI mockup\", \"subject\": { \"description\": \"photorealistic young Asian woman, sweaty glowing skin, long dark w",
    "promptPreview": "{ \"type\": \"e-commerce livestream UI mockup\", \"subject\": { \"description\": \"photorealistic young Asian woman, sweaty glowing skin, long dark wavy hair, wearing a ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@onlyhuman028",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Realistic",
      "Product",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-158",
    "locale": "zh",
    "title": "界面交互设计图",
    "category": "UI 与界面",
    "description": "{ \"type\": \"e-commerce live stream interface mockup\", \"subject\": { \"description\": \"young Asian woman, long wavy dark hair, wearing a white sh",
    "promptPreview": "{ \"type\": \"e-commerce live stream interface mockup\", \"subject\": { \"description\": \"young Asian woman, long wavy dark hair, wearing a white short-sleeve polo shir",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@coconut_256",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Realistic",
      "Product",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-133",
    "locale": "zh",
    "title": "界面交互设计图",
    "category": "UI 与界面",
    "description": "{ \"type\": \"brand identity system presentation board\", \"header\": { \"title\": \"品牌视觉识别系统 BRAND IDENTITY SYSTEM\", \"slogan\": \"爱它·懂它·陪伴它\" }, \"main_",
    "promptPreview": "{ \"type\": \"brand identity system presentation board\", \"header\": { \"title\": \"品牌视觉识别系统 BRAND IDENTITY SYSTEM\", \"slogan\": \"爱它·懂它·陪伴它\" }, \"main_logo\": { \"text\": \"{a",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@yyyole",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-131",
    "locale": "en",
    "title": "界面交互设计图",
    "category": "UI 与界面",
    "description": "{ \"type\": \"UI/UX landing page mockup\", \"theme\": \"dark mode, sleek modern aesthetic, glassmorphism, {argument name=\\\"primary accent color\\\" d",
    "promptPreview": "{ \"type\": \"UI/UX landing page mockup\", \"theme\": \"dark mode, sleek modern aesthetic, glassmorphism, {argument name=\\\"primary accent color\\\" default=\\\"neon purple",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@IndieDevHailey",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Brand",
      "3D",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-104",
    "locale": "en",
    "title": "界面交互设计图",
    "category": "UI 与界面",
    "description": "{ \"type\": \"YouTube livestream UI\", \"top_nav\": { \"logo\": \"YouTube Premium\", \"search\": \"{argument name=\\\"search query\\\" default=\\\"bilal fraiha",
    "promptPreview": "{ \"type\": \"YouTube livestream UI\", \"top_nav\": { \"logo\": \"YouTube Premium\", \"search\": \"{argument name=\\\"search query\\\" default=\\\"bilal fraiha\\\"}\", \"icons\": 3 }, ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@marouane53",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Brand",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-57",
    "locale": "en",
    "title": "界面交互设计图",
    "category": "UI 与界面",
    "description": "{ \"type\": \"mobile social media app UI mockup\", \"platform\": \"Twitter/X dark mode\", \"header\": { \"status_bar\": \"time 19:28, bird icon, signal, ",
    "promptPreview": "{ \"type\": \"mobile social media app UI mockup\", \"platform\": \"Twitter/X dark mode\", \"header\": { \"status_bar\": \"time 19:28, bird icon, signal, wifi, battery\", \"nav",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Character",
      "Classical",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-289",
    "locale": "zh",
    "title": "直播界面设计图",
    "category": "UI 与界面",
    "description": "[中文] 生成特朗普和金正恩在抖音直播间打PK的截图 [English] Generate a screenshot of Trump and Kim Jong-un doing a PK battle in a TikTok live stream room",
    "promptPreview": "[中文] 生成特朗普和金正恩在抖音直播间打PK的截图 [English] Generate a screenshot of Trump and Kim Jong-un doing a PK battle in a TikTok live stream room",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@rionaifantasy",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-152",
    "locale": "zh",
    "title": "直播界面设计图",
    "category": "UI 与界面",
    "description": "{ \"type\": \"e-commerce livestream screenshot mockup\", \"scene\": { \"subject\": \"{argument name=\\\"main subject\\\" default=\\\"Caucasian male resembl",
    "promptPreview": "{ \"type\": \"e-commerce livestream screenshot mockup\", \"scene\": { \"subject\": \"{argument name=\\\"main subject\\\" default=\\\"Caucasian male resembling Sam Altman\\\"}\", ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@coder_left",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Product",
      "Brand",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-149",
    "locale": "zh",
    "title": "直播界面设计图",
    "category": "UI 与界面",
    "description": "{ \"type\": \"mobile livestream e-commerce interface mockup\", \"subject\": { \"person\": \"Elon Musk\", \"clothing\": \"black t-shirt with SPACEX logo\",",
    "promptPreview": "{ \"type\": \"mobile livestream e-commerce interface mockup\", \"subject\": { \"person\": \"Elon Musk\", \"clothing\": \"black t-shirt with SPACEX logo\", \"pose\": \"gesturing ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@JCutcut47692",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Realistic",
      "Product",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-49",
    "locale": "en",
    "title": "直播界面设计图",
    "category": "UI 与界面",
    "description": "A 9:16 aspect ratio image, generating a screenshot of a Douyin livestream where {argument name=\"celebrity\" default=\"Liu Yifei\"} is broadcast",
    "promptPreview": "A 9:16 aspect ratio image, generating a screenshot of a Douyin livestream where {argument name=\"celebrity\" default=\"Liu Yifei\"} is broadcasting, holding a sign ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@kylegeeks",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-2"
  },
  {
    "id": "freestylefly-285",
    "locale": "zh",
    "title": "真实动漫画面快照",
    "category": "UI 与界面",
    "description": "[中文] 向我展示这张附带的图像作为一部真实动漫的快照 [English] Show me the attached image as a snapshot from an actual anime",
    "promptPreview": "[中文] 向我展示这张附带的图像作为一部真实动漫的快照 [English] Show me the attached image as a snapshot from an actual anime",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@Thereallo1026",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-260",
    "locale": "zh",
    "title": "社媒界面截图",
    "category": "UI 与界面",
    "description": "[中文] 生成抖音内容截图，主题：跟上AI浪潮9.9包教会，iPhone尺寸 [English] Generate a screenshot of Douyin content, theme: Catch up with the AI wave, 9.9 to learn it ",
    "promptPreview": "[中文] 生成抖音内容截图，主题：跟上AI浪潮9.9包教会，iPhone尺寸 [English] Generate a screenshot of Douyin content, theme: Catch up with the AI wave, 9.9 to learn it all, iPhone size",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Tech",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-111",
    "locale": "en",
    "title": "视频封面界面图",
    "category": "UI 与界面",
    "description": "A YouTube thumbnail-style collage for a {argument name=\"overall mood\" default=\"dark, dramatic, true crime investigation\"}. In the center is ",
    "promptPreview": "A YouTube thumbnail-style collage for a {argument name=\"overall mood\" default=\"dark, dramatic, true crime investigation\"}. In the center is a highly detailed, c",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@mirochill",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Fashion",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-110",
    "locale": "en",
    "title": "视频封面界面图",
    "category": "UI 与界面",
    "description": "Thumbnail for a YouTube unboxing video, a video of {argument name=\"topic\" default=\"opening all overdue bills\"}, {argument name=\"quantity\" de",
    "promptPreview": "Thumbnail for a YouTube unboxing video, a video of {argument name=\"topic\" default=\"opening all overdue bills\"}, {argument name=\"quantity\" default=\"100 in a row\"",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@TlanoAI",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Poster",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-103",
    "locale": "en",
    "title": "视频封面界面图",
    "category": "UI 与界面",
    "description": "{argument name=\"pianist\" default=\"Vladimir Horowitz\"} performs a {argument name=\"event\" default=\"live piano recital\"} streamed on {argument ",
    "promptPreview": "{argument name=\"pianist\" default=\"Vladimir Horowitz\"} performs a {argument name=\"event\" default=\"live piano recital\"} streamed on {argument name=\"platform\" defa",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@bowowwoaaa2",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Poster",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-92",
    "locale": "en",
    "title": "视频封面界面图",
    "category": "UI 与界面",
    "description": "An anime-style YouTube stream thumbnail featuring a cheerful female VTuber. She has long {argument name=\"hair color\" default=\"pink with ligh",
    "promptPreview": "An anime-style YouTube stream thumbnail featuring a cheerful female VTuber. She has long {argument name=\"hair color\" default=\"pink with light blue inner highlig",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Yuupapa_free",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Poster",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-209",
    "locale": "zh",
    "title": "神话三国枪战世界",
    "category": "UI 与界面",
    "description": "[中文] 模仿《无畏契约》的风格，生成一个三国神话的 FPS 游戏 [English] Imitating the style of Valorant, generate a Three Kingdoms mythological FPS game",
    "promptPreview": "[中文] 模仿《无畏契约》的风格，生成一个三国神话的 FPS 游戏 [English] Imitating the style of Valorant, generate a Three Kingdoms mythological FPS game",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@op7418",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-259",
    "locale": "zh",
    "title": "精致女孩背后的网贷真相",
    "category": "UI 与界面",
    "description": "[中文] 生成小红书内容截图，主题：精致女孩背后都有网贷，iPhone尺寸 [English] Generate Xiaohongshu content screenshot, theme: Behind every exquisite girl there is online ",
    "promptPreview": "[中文] 生成小红书内容截图，主题：精致女孩背后都有网贷，iPhone尺寸 [English] Generate Xiaohongshu content screenshot, theme: Behind every exquisite girl there is online loan, iPhone size",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-94",
    "locale": "en",
    "title": "绘画艺术风格图",
    "category": "UI 与界面",
    "description": "{ \"type\": \"VTuber stream thumbnail\", \"theme\": \"pastel pink, soft, cute, lace, ribbons, hearts, bunny motif\", \"character\": { \"position\": \"rig",
    "promptPreview": "{ \"type\": \"VTuber stream thumbnail\", \"theme\": \"pastel pink, soft, cute, lace, ribbons, hearts, bunny motif\", \"character\": { \"position\": \"right side, waist-up\", ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@sayaka_aiart",
    "tags": [
      "UI & Interfaces",
      "Illustration",
      "Character",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-148",
    "locale": "en",
    "title": "综合应用场景图",
    "category": "UI 与界面",
    "description": "A {argument name=\"platform\" default=\"Taobao\"} product detail page for {argument name=\"robot model\" default=\"T-800 robot\"}, displaying: front",
    "promptPreview": "A {argument name=\"platform\" default=\"Taobao\"} product detail page for {argument name=\"robot model\" default=\"T-800 robot\"}, displaying: front, side, and back thr",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@alanlovelq",
    "tags": [
      "UI & Interfaces",
      "Product",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-109",
    "locale": "en",
    "title": "综合应用场景图",
    "category": "UI 与界面",
    "description": "{argument name=\"subject\" default=\"A beautiful internet celebrity\"} is live-streaming a {argument name=\"activity\" default=\"game\"}.",
    "promptPreview": "{argument name=\"subject\" default=\"A beautiful internet celebrity\"} is live-streaming a {argument name=\"activity\" default=\"game\"}.",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@underwoodxie96",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-249",
    "locale": "zh",
    "title": "美女举牌感谢大哥打赏大火箭",
    "category": "UI 与界面",
    "description": "[中文] 生成一个抖音直播的截图 ，一个美女在直播，美女手里拿着牌子，上面写着：谢谢行者大哥的大火箭！ [English] Generate a screenshot of a TikTok live stream, a beautiful woman is live strea",
    "promptPreview": "[中文] 生成一个抖音直播的截图 ，一个美女在直播，美女手里拿着牌子，上面写着：谢谢行者大哥的大火箭！ [English] Generate a screenshot of a TikTok live stream, a beautiful woman is live streaming, the beautiful ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@joshesye",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-4",
    "locale": "zh",
    "title": "老干妈风味",
    "category": "UI 与界面",
    "description": "特朗普在抖音直播间卖老干妈，手里举着「老干妈风味」新品，背景还是 SpaceX 那种科技感，左下角弹幕飘着「特斯拉车主：求上链接」。",
    "promptPreview": "特朗普在抖音直播间卖老干妈，手里举着「老干妈风味」新品，背景还是 SpaceX 那种科技感，左下角弹幕飘着「特斯拉车主：求上链接」。",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "小红书号989137706",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-197",
    "locale": "zh",
    "title": "英雄联盟特朗普中路对决哈梅内伊",
    "category": "UI 与界面",
    "description": "[中文] 帮我生成一张特朗普对战哈梅内伊在英雄联盟中路对线的截图。 [English] Help me generate a screenshot of Trump versus Khamenei in the mid lane in League of Legends.",
    "promptPreview": "[中文] 帮我生成一张特朗普对战哈梅内伊在英雄联盟中路对线的截图。 [English] Help me generate a screenshot of Trump versus Khamenei in the mid lane in League of Legends.",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@underwoodxie96",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-163",
    "locale": "zh",
    "title": "诗仙李白月下直播起舞",
    "category": "UI 与界面",
    "description": "[中文] 李白在抖音直播月下起舞 [English] Li Bai dancing under the moon during a Douyin livestream",
    "promptPreview": "[中文] 李白在抖音直播月下起舞 [English] Li Bai dancing under the moon during a Douyin livestream",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-404",
    "locale": "zh",
    "title": "豪华社媒破屏商业广告",
    "category": "UI 与界面",
    "description": "动态的豪华商业广告海报，特色是超现实3D渲染的充满活力的年轻女性，以上传的女性面部作为参考，穿着高级亮橙色设计师服装、豪华配饰以及时尚的金色墨镜，自信地从一个巨大的金色智能手机屏幕中爆发出色。她的姿势有力且时尚，一只运动鞋通过强烈的强制透视戏剧性地穿过数字显示屏，踏入现实。 构图",
    "promptPreview": "动态的豪华商业广告海报，特色是超现实3D渲染的充满活力的年轻女性，以上传的女性面部作为参考，穿着高级亮橙色设计师服装、豪华配饰以及时尚的金色墨镜，自信地从一个巨大的金色智能手机屏幕中爆发出色。她的姿势有力且时尚，一只运动鞋通过强烈的强制透视戏剧性地穿过数字显示屏，踏入现实。 构图在前台强调她白色豪华运动鞋，带有纹理口香",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@you1873118",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Poster",
      "Brand",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-213",
    "locale": "zh",
    "title": "金瓶梅古风开放世界游戏截图",
    "category": "UI 与界面",
    "description": "[中文] 帮我生成一个以《金瓶梅》为主题的古代 ARPG MMO 开放世界游戏的截图 [English] Help me generate a screenshot of an ancient ARPG MMO open-world game themed around Jin ",
    "promptPreview": "[中文] 帮我生成一个以《金瓶梅》为主题的古代 ARPG MMO 开放世界游戏的截图 [English] Help me generate a screenshot of an ancient ARPG MMO open-world game themed around Jin Ping Mei.",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@op7418",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Classical",
      "Story",
      "History"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-437",
    "locale": "en",
    "title": "面部美学分析报告",
    "category": "UI 与界面",
    "description": "Create a clean, minimal, luxury-style facial aesthetics analysis report based on the uploaded portrait photo. Design style: Ultra-modern bla",
    "promptPreview": "Create a clean, minimal, luxury-style facial aesthetics analysis report based on the uploaded portrait photo. Design style: Ultra-modern black and white interfa",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@saniaspeaks_",
    "tags": [
      "UI & Interfaces",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-207",
    "locale": "zh",
    "title": "黑神话潘金莲绝美游戏封面",
    "category": "UI 与界面",
    "description": "[中文] 生成一张黑神话·潘金莲的游戏介绍画面，人物十分的迷人 [English] Generate a game introduction screen for Black Myth: Pan Jinlian, the character is extremely charmi",
    "promptPreview": "[中文] 生成一张黑神话·潘金莲的游戏介绍画面，人物十分的迷人 [English] Generate a game introduction screen for Black Myth: Pan Jinlian, the character is extremely charming.",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "UI & Interfaces",
      "Poster",
      "Character",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-302",
    "locale": "zh",
    "title": "九位大师的机械键盘设计图鉴",
    "category": "其他案例",
    "description": "[中文] 一个九宫格图片，展现九位当代知名设计师设计的同一组物体：机械键盘，包括设计师头像，设计师对于设计的中文文字解读和作品呈现。排版统一规则 [English] A nine-grid image showing the same group of objects desig",
    "promptPreview": "[中文] 一个九宫格图片，展现九位当代知名设计师设计的同一组物体：机械键盘，包括设计师头像，设计师对于设计的中文文字解读和作品呈现。排版统一规则 [English] A nine-grid image showing the same group of objects designed by nine contempo",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@TanShilong",
    "tags": [
      "Other Use Cases",
      "Character",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-252",
    "locale": "zh",
    "title": "五一劳动节手举牌创意设计集",
    "category": "其他案例",
    "description": "[中文] 生成一系列五一劳动节的手举牌设计 [English] Generate a series of hand-held sign designs for May Day Labor Day",
    "promptPreview": "[中文] 生成一系列五一劳动节的手举牌设计 [English] Generate a series of hand-held sign designs for May Day Labor Day",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@akokoi1",
    "tags": [
      "Other Use Cases",
      "Other Use Cases",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-241",
    "locale": "zh",
    "title": "关键人物关系图谱",
    "category": "其他案例",
    "description": "[中文] 请你生成 《XXX》 的关键人物关系图。 [English] Please generate a key character relationship diagram for \"XXX\".",
    "promptPreview": "[中文] 请你生成 《XXX》 的关键人物关系图。 [English] Please generate a key character relationship diagram for \"XXX\".",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@yihui_indie",
    "tags": [
      "Other Use Cases",
      "Infographic",
      "Character",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-81",
    "locale": "en",
    "title": "写实摄影风格图",
    "category": "其他案例",
    "description": "{ \"type\": \"scientific hardware diagram\", \"layout\": { \"main_scene\": \"3D render of an optical table with a red laser beam passing through 11 a",
    "promptPreview": "{ \"type\": \"scientific hardware diagram\", \"layout\": { \"main_scene\": \"3D render of an optical table with a red laser beam passing through 11 aligned optical compo",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@HumanOS_v2",
    "tags": [
      "Other Use Cases",
      "Infographic",
      "Realistic",
      "3D",
      "Tech",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-309",
    "locale": "zh",
    "title": "创意树叶拼贴构成的角色画像",
    "category": "其他案例",
    "description": "[中文] { 角色名称 } 完全由天然树叶制成，创意树叶拼贴艺术，分层绿叶和干叶构成身体、面部和衣服，可见叶脉和纹理，手工植物艺术风格，干净的白色背景，俯视平铺构图，高度细节，柔和自然光，逼真树叶纹理，8k [English] { CHARACTER NAME } made en",
    "promptPreview": "[中文] { 角色名称 } 完全由天然树叶制成，创意树叶拼贴艺术，分层绿叶和干叶构成身体、面部和衣服，可见叶脉和纹理，手工植物艺术风格，干净的白色背景，俯视平铺构图，高度细节，柔和自然光，逼真树叶纹理，8k [English] { CHARACTER NAME } made entirely from natural ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@meng_dagg695",
    "tags": [
      "Other Use Cases",
      "Realistic",
      "Character",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-368",
    "locale": "zh",
    "title": "印度餐厅菜单改造宣传图",
    "category": "其他案例",
    "description": "这是india 料理中的一份真实menu。根据此 重新生成带文本说明的 引人入胜垂涎欲滴的 说明图片 先用English 文本易于识别（手机小屏幕） 这个是beef roast",
    "promptPreview": "这是india 料理中的一份真实menu。根据此 重新生成带文本说明的 引人入胜垂涎欲滴的 说明图片 先用English 文本易于识别（手机小屏幕） 这个是beef roast",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@Johnson998877",
    "tags": [
      "Other Use Cases",
      "Other Use Cases",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-290",
    "locale": "zh",
    "title": "古风诗人镭射典藏卡牌",
    "category": "其他案例",
    "description": "[中文] 为中国古代诗人设计一套游戏卡片，并按照SSR SR R 分级，重点卡片有放大展示的效果，包括卡面设计和人物介绍，有很高级的游戏卡片质感，稀有卡片还会有特色的光影例如镭射效果 需要有套卡设计和技能设计，并附带较为详细的说明 [English] Design a set o",
    "promptPreview": "[中文] 为中国古代诗人设计一套游戏卡片，并按照SSR SR R 分级，重点卡片有放大展示的效果，包括卡面设计和人物介绍，有很高级的游戏卡片质感，稀有卡片还会有特色的光影例如镭射效果 需要有套卡设计和技能设计，并附带较为详细的说明 [English] Design a set of game cards for anc",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@TanShilong",
    "tags": [
      "Other Use Cases",
      "UI",
      "Character",
      "Classical",
      "Tech",
      "Commerce",
      "History"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-391",
    "locale": "en",
    "title": "四国文化锚点服装设计宫格",
    "category": "其他案例",
    "description": "<instructions> input: continent pick 4 lesser known countries in that continent function drawx($lesser known country){ > Anchor 1: \"$lesser ",
    "promptPreview": "<instructions> input: continent pick 4 lesser known countries in that continent function drawx($lesser known country){ > Anchor 1: \"$lesser known's famous archi",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Gdgtify",
    "tags": [
      "Other Use Cases",
      "3D",
      "Fashion",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-80",
    "locale": "en",
    "title": "图像生成案例图",
    "category": "其他案例",
    "description": "[CORE TASK] Transform the provided input image into a pose-and-light analysis sheet. This is NOT a finished character illustration. This is ",
    "promptPreview": "[CORE TASK] Transform the provided input image into a pose-and-light analysis sheet. This is NOT a finished character illustration. This is NOT a clothing sheet",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@WOZ1Tx2JZ3kCeBj",
    "tags": [
      "Other Use Cases",
      "UI",
      "Realistic",
      "Illustration",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-295",
    "locale": "zh",
    "title": "复古传统老黄历二零二六年四月十八",
    "category": "其他案例",
    "description": "[中文] 生成一张2026年4月18日的老黄历 [English] Generate an old almanac for April 18, 2026",
    "promptPreview": "[中文] 生成一张2026年4月18日的老黄历 [English] Generate an old almanac for April 18, 2026",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "Other Use Cases",
      "Other Use Cases",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-250",
    "locale": "zh",
    "title": "小王子与星舰的浪漫联名",
    "category": "其他案例",
    "description": "[中文] 设计一张小王子和SpaceX联名的明信片 [English] Design a postcard co-branded by The Little Prince and SpaceX",
    "promptPreview": "[中文] 设计一张小王子和SpaceX联名的明信片 [English] Design a postcard co-branded by The Little Prince and SpaceX",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@akokoi1",
    "tags": [
      "Other Use Cases",
      "Brand",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-323",
    "locale": "en",
    "title": "应用界面样机图",
    "category": "其他案例",
    "description": "Create a hyper-realistic, cinematic Instagram post layout where the Instagram UI exists as a physical, tangible 3D object, photographed like",
    "promptPreview": "Create a hyper-realistic, cinematic Instagram post layout where the Instagram UI exists as a physical, tangible 3D object, photographed like a premium commercia",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Mystveil7",
    "tags": [
      "Other Use Cases",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-106",
    "locale": "en",
    "title": "应用界面样机图",
    "category": "其他案例",
    "description": "{ \"type\": \"YouTube thumbnail graphic\", \"style\": \"anime, edgy, neon pink and black color scheme, grunge and splatter accents\", \"character\": {",
    "promptPreview": "{ \"type\": \"YouTube thumbnail graphic\", \"style\": \"anime, edgy, neon pink and black color scheme, grunge and splatter accents\", \"character\": { \"appearance\": \"anim",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@abdiisan",
    "tags": [
      "Other Use Cases",
      "UI",
      "Poster",
      "Brand",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-297",
    "locale": "zh",
    "title": "手写食谱变身杂志级跨页",
    "category": "其他案例",
    "description": "[中文] 手写食谱 → 专业食谱页面 上传一份凌乱的手写家庭食谱；模型会搜索准确的现代计量/营养信息，然后生成一份精致的杂志风格双页跨页，包含分步平铺图、完美的食材标签和卡路里分解。 [INSERT_RECIPE_LINK] [English] Handwritten Recip",
    "promptPreview": "[中文] 手写食谱 → 专业食谱页面 上传一份凌乱的手写家庭食谱；模型会搜索准确的现代计量/营养信息，然后生成一份精致的杂志风格双页跨页，包含分步平铺图、完美的食材标签和卡路里分解。 [INSERT_RECIPE_LINK] [English] Handwritten Recipe → Professional Coo",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@maxescu",
    "tags": [
      "Other Use Cases",
      "Other Use Cases",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-204",
    "locale": "zh",
    "title": "智能动画分镜生成器",
    "category": "其他案例",
    "description": "[中文] 生成一张动画分镜生成器 [English] Generate an animation storyboard generator",
    "promptPreview": "[中文] 生成一张动画分镜生成器 [English] Generate an animation storyboard generator",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@joshesye",
    "tags": [
      "Other Use Cases",
      "Other Use Cases",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-203",
    "locale": "zh",
    "title": "杠精视角的独特文案创意",
    "category": "其他案例",
    "description": "[中文] 杠精视角文案 + GPT Image 2 [English] Troll perspective copywriting + GPT Image 2",
    "promptPreview": "[中文] 杠精视角文案 + GPT Image 2 [English] Troll perspective copywriting + GPT Image 2",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@joshesye",
    "tags": [
      "Other Use Cases",
      "Other Use Cases",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-286",
    "locale": "zh",
    "title": "珠江新城剪纸璀璨夜景",
    "category": "其他案例",
    "description": "[中文] 以珠江新城现代都市景观为灵感的剪纸艺术，通过精巧的镂空手法在一整幅纸上，立体刻画广州塔、东西双塔等地标建筑与繁华城景。 所有建筑与元素均以流畅的线条与结构相连，无孤立部分，构成一幅完整的都市画卷。 画面采用金属箔或光泽纸材质，表面带有细腻的明暗光泽，在光照下呈现柔和的高",
    "promptPreview": "[中文] 以珠江新城现代都市景观为灵感的剪纸艺术，通过精巧的镂空手法在一整幅纸上，立体刻画广州塔、东西双塔等地标建筑与繁华城景。 所有建筑与元素均以流畅的线条与结构相连，无孤立部分，构成一幅完整的都市画卷。 画面采用金属箔或光泽纸材质，表面带有细腻的明暗光泽，在光照下呈现柔和的高光与阴影，仿佛被城市灯光轻轻照亮。 背景",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "Other Use Cases",
      "UI",
      "Realistic",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-363",
    "locale": "en",
    "title": "磁场铁粉 Logo 物理成像",
    "category": "其他案例",
    "description": "Transform the uploaded logo into a hyper-realistic scene where the logo silhouette is formed by iron filings reacting to a magnetic field. T",
    "promptPreview": "Transform the uploaded logo into a hyper-realistic scene where the logo silhouette is formed by iron filings reacting to a magnetic field. The logo must keep it",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Naiknelofar788",
    "tags": [
      "Other Use Cases",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-242",
    "locale": "zh",
    "title": "绝美国风工笔画书签设计",
    "category": "其他案例",
    "description": "[中文] 生成一系列工笔画书签的设计稿 [English] Generate a series of design drafts for Gongbi painting bookmarks.",
    "promptPreview": "[中文] 生成一系列工笔画书签的设计稿 [English] Generate a series of design drafts for Gongbi painting bookmarks.",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@akokoi1",
    "tags": [
      "Other Use Cases",
      "Illustration",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-215",
    "locale": "zh",
    "title": "西方艺术演进像素博物馆",
    "category": "其他案例",
    "description": "[中文] 创作一张超高细节等距像素艺术时间线插画（3:4，4K），融合细节密度、象征性与隐喻。用户指定的主题为【Western Art Development】。 首先，围绕Western Art Development进行推理，确定：主题的中英文标题、涵盖的最早与最近历史时期、",
    "promptPreview": "[中文] 创作一张超高细节等距像素艺术时间线插画（3:4，4K），融合细节密度、象征性与隐喻。用户指定的主题为【Western Art Development】。 首先，围绕Western Art Development进行推理，确定：主题的中英文标题、涵盖的最早与最近历史时期、起始阶段标签与结束阶段标签，以及3-5个",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@GeekCatX",
    "tags": [
      "Other Use Cases",
      "UI",
      "Poster",
      "Illustration",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-251",
    "locale": "zh",
    "title": "言叶之庭春雨绿意单日历",
    "category": "其他案例",
    "description": "[中文] 生成一张言叶之庭2026年4月19日单日日历 [English] Generate a single-day calendar for The Garden of Words on April 19, 2026",
    "promptPreview": "[中文] 生成一张言叶之庭2026年4月19日单日日历 [English] Generate a single-day calendar for The Garden of Words on April 19, 2026",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@akokoi1",
    "tags": [
      "Other Use Cases",
      "Other Use Cases",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-196",
    "locale": "zh",
    "title": "试卷上的涂鸦巨龙",
    "category": "其他案例",
    "description": "[中文] 一个巨大的巨龙，庞大的规模，高耸的存在感， 一个远超人类尺寸的巨大实体，压倒性和压迫性的， 用极其密集的混乱涂鸦线条绘制， 超密集的重叠笔触，纠缠和混乱的线条画， 在真实的印刷英文/中文教科书或试卷页面上， 可见的文本、布局和纸张纹理清晰透出， 圆珠笔绘画风格，精细的墨",
    "promptPreview": "[中文] 一个巨大的巨龙，庞大的规模，高耸的存在感， 一个远超人类尺寸的巨大实体，压倒性和压迫性的， 用极其密集的混乱涂鸦线条绘制， 超密集的重叠笔触，纠缠和混乱的线条画， 在真实的印刷英文/中文教科书或试卷页面上， 可见的文本、布局和纸张纹理清晰透出， 圆珠笔绘画风格，精细的墨水线条，杂乱的分层笔触， 没有干净的轮廓",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@GeekCatX",
    "tags": [
      "Other Use Cases",
      "Illustration",
      "Tech",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-173",
    "locale": "zh",
    "title": "银河繁星点缀的冰蓝襦裙",
    "category": "其他案例",
    "description": "[中文] 服裝細節： 模特兒身穿一套精緻的淡冰藍色齊胸襦裙，採用多層輕盈的薄紗和絲綢歐根紗材質制成。其寬大的、半透明的廣袖上點綴著如繁星般微小的銀色和淺藍色亮片刺繡，在光線下閃爍（具有銀河般的夢幻感）。抹胸位置有複雜的銀色蕾絲和編織紋理細節，腰帶自然垂落。 材質與光影： 畫面呈現",
    "promptPreview": "[中文] 服裝細節： 模特兒身穿一套精緻的淡冰藍色齊胸襦裙，採用多層輕盈的薄紗和絲綢歐根紗材質制成。其寬大的、半透明的廣袖上點綴著如繁星般微小的銀色和淺藍色亮片刺繡，在光線下閃爍（具有銀河般的夢幻感）。抹胸位置有複雜的銀色蕾絲和編織紋理細節，腰帶自然垂落。 材質與光影： 畫面呈現 8k 超高分辨率和對織物微距紋理的極致",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@fdtreesky",
    "tags": [
      "Other Use Cases",
      "UI",
      "3D",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-385",
    "locale": "zh",
    "title": "青岛啤酒灵感女装系列",
    "category": "其他案例",
    "description": "Inspired by Tsingtao (China beer)🍺 “Inspired by this product, design a set of cool-style women's clothing”",
    "promptPreview": "Inspired by Tsingtao (China beer)🍺 “Inspired by this product, design a set of cool-style women's clothing”",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Popcraft_ai",
    "tags": [
      "Other Use Cases",
      "Product",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-245",
    "locale": "zh",
    "title": "马斯克专属篆刻印章设计",
    "category": "其他案例",
    "description": "[中文] 给”埃隆·马斯克”设计一组篆刻印章 [English] Design a set of seal carving stamps for \"Elon Musk\"",
    "promptPreview": "[中文] 给”埃隆·马斯克”设计一组篆刻印章 [English] Design a set of seal carving stamps for \"Elon Musk\"",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@akokoi1",
    "tags": [
      "Other Use Cases",
      "Other Use Cases",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-337",
    "locale": "zh",
    "title": "《短歌行》诗词意境图",
    "category": "历史与古典主题",
    "description": "帮我生成一张《短歌行》的意境图，带整篇《短歌行》文字",
    "promptPreview": "帮我生成一张《短歌行》的意境图，带整篇《短歌行》文字",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "苍何原创实测（公众号文章《我逆向了 329 条 GPT-Image2 提示词模板，全部开源！》）",
    "tags": [
      "History & Classical Themes",
      "History",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-415",
    "locale": "zh",
    "title": "东方神话人物志百科海报",
    "category": "历史与古典主题",
    "description": "请基于【孙悟空】生成一张竖版 A4「东方神话人物志」百科海报。 你是一名东方神话视觉设定师、古籍图谱设计师和中文信息图设计师。请根据【人物名称】在东方神话、民间传说、古籍、戏曲、文学或传统文化中的已知形象，生成一张内容完整、视觉统一、具有东方神话气质的人物档案海报。 最重要原则：",
    "promptPreview": "请基于【孙悟空】生成一张竖版 A4「东方神话人物志」百科海报。 你是一名东方神话视觉设定师、古籍图谱设计师和中文信息图设计师。请根据【人物名称】在东方神话、民间传说、古籍、戏曲、文学或传统文化中的已知形象，生成一张内容完整、视觉统一、具有东方神话气质的人物档案海报。 最重要原则： 全部内容必须围绕【人物名称】展开，不要",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@TanLuAI",
    "tags": [
      "History & Classical Themes",
      "UI",
      "Infographic",
      "Poster",
      "Fashion",
      "Story",
      "History"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-44",
    "locale": "en",
    "title": "古风历史题材图",
    "category": "历史与古典主题",
    "description": "Generate avatars of various emperors from the {argument name=\"dynasty\" default=\"Ming Dynasty\"} based on the style of the uploaded image, wit",
    "promptPreview": "Generate avatars of various emperors from the {argument name=\"dynasty\" default=\"Ming Dynasty\"} based on the style of the uploaded image, with their posthumous n",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "History & Classical Themes",
      "Character",
      "Classical",
      "Commerce",
      "History"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-226",
    "locale": "zh",
    "title": "古风明朝帝王群像长卷",
    "category": "历史与古典主题",
    "description": "[中文] 根据上传图片的风格，生成明朝各个皇帝的头像，头像下面有他们的谥号和名字 [English] Based on the style of the uploaded image, generate portraits of the emperors of the Ming ",
    "promptPreview": "[中文] 根据上传图片的风格，生成明朝各个皇帝的头像，头像下面有他们的谥号和名字 [English] Based on the style of the uploaded image, generate portraits of the emperors of the Ming Dynasty, with their ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "History & Classical Themes",
      "Classical",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-174",
    "locale": "zh",
    "title": "唐朝贵妇遛粉色马甲异形工笔画",
    "category": "历史与古典主题",
    "description": "[中文] 一幅细节丰富的工笔画，描绘了一位唐朝贵族女子在御花园中漫步。她看起来优雅而平静。 她手里拿着一根金色的牵引绳。牵引绳的尽头是一只可怕的**异形怪物（出自电影《异形》）**。然而，这只异形穿着一件**可爱的粉色丝绸马甲**，并且表现得像一只训练有素的狗。 背景有牡丹和蝴蝶",
    "promptPreview": "[中文] 一幅细节丰富的工笔画，描绘了一位唐朝贵族女子在御花园中漫步。她看起来优雅而平静。 她手里拿着一根金色的牵引绳。牵引绳的尽头是一只可怕的**异形怪物（出自电影《异形》）**。然而，这只异形穿着一件**可爱的粉色丝绸马甲**，并且表现得像一只训练有素的狗。 背景有牡丹和蝴蝶。 **在右下角，有一个红色的竖排艺术家",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@johnAGI168",
    "tags": [
      "History & Classical Themes",
      "Illustration",
      "Classical",
      "Tech",
      "Commerce",
      "History"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-206",
    "locale": "zh",
    "title": "国风工笔八仙长卷插画",
    "category": "历史与古典主题",
    "description": "[中文] （国风卷轴插画师）你是一位顶尖的中国传统工笔人物画师，擅长将经典人物群像绘制成长卷式百科海报。根据用户指定的【eight immortals】，生成一张 “中国传统人物群像长卷海报”：画面为横向长卷式构图，所有人物排成一条队列，从左至右依次展开；每个人物都有鲜明的传统服",
    "promptPreview": "[中文] （国风卷轴插画师）你是一位顶尖的中国传统工笔人物画师，擅长将经典人物群像绘制成长卷式百科海报。根据用户指定的【eight immortals】，生成一张 “中国传统人物群像长卷海报”：画面为横向长卷式构图，所有人物排成一条队列，从左至右依次展开；每个人物都有鲜明的传统服饰、标志性道具和神态，下方配有竖排名牌标",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@GeekCatX",
    "tags": [
      "History & Classical Themes",
      "Poster",
      "Illustration",
      "Brand",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-167",
    "locale": "zh",
    "title": "大唐玄武门之变的朋友圈",
    "category": "历史与古典主题",
    "description": "[中文] 玄武门之变的朋友圈 [English] WeChat Moments of the Xuanwu Gate Incident",
    "promptPreview": "[中文] 玄武门之变的朋友圈 [English] WeChat Moments of the Xuanwu Gate Incident",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@Tz_2022",
    "tags": [
      "History & Classical Themes",
      "History",
      "Social",
      "History"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-268",
    "locale": "zh",
    "title": "威化岛回军前夕李成桂动态",
    "category": "历史与古典主题",
    "description": "[中文] 태조 이성계의 X 페이지(위화도 회군을 벌이기 직전- 최영 장군과 서로 디스하는 내용이 담긴 게시글들)을 만들어 주세요. [English] Please create an X page of King Taejo Yi Seong-gye (right",
    "promptPreview": "[中文] 태조 이성계의 X 페이지(위화도 회군을 벌이기 직전- 최영 장군과 서로 디스하는 내용이 담긴 게시글들)을 만들어 주세요. [English] Please create an X page of King Taejo Yi Seong-gye (right before carrying out",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@SKA_Neotype",
    "tags": [
      "History & Classical Themes",
      "History",
      "Tech",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-267",
    "locale": "zh",
    "title": "宋朝文人的赛博朋友圈",
    "category": "历史与古典主题",
    "description": "[中文] \"宋朝人的朋友圈\"/\"SONG DYNASTY SOCIAL MEDIA FEED\"，古今穿越幽默融合界面设计风格，画面模拟手机社交媒体界面，但内容全部是宋朝场景头像是宋代文人画像，用户名\"苏东坡SuShi_Official\"，发布内容\"刚到黄州，被贬了但心情还行。今天",
    "promptPreview": "[中文] \"宋朝人的朋友圈\"/\"SONG DYNASTY SOCIAL MEDIA FEED\"，古今穿越幽默融合界面设计风格，画面模拟手机社交媒体界面，但内容全部是宋朝场景头像是宋代文人画像，用户名\"苏东坡SuShi_Official\"，发布内容\"刚到黄州，被贬了但心情还行。今天自己做了东坡肉，味道绝了，附菜谱：\"，配",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@Panda20230902",
    "tags": [
      "History & Classical Themes",
      "UI",
      "Illustration",
      "Character",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-292",
    "locale": "zh",
    "title": "明朝登基宝玉的推文页面",
    "category": "历史与古典主题",
    "description": "[中文] 创建一个宝玉（查阅 https://x.com/dotey 这个推主的主页及部分推文）穿越到明朝，登基之后依据其业务/个性，绘制的其新的X帖子页面。 [English] Create a new X post page illustrated for Baoyu (re",
    "promptPreview": "[中文] 创建一个宝玉（查阅 https://x.com/dotey 这个推主的主页及部分推文）穿越到明朝，登基之后依据其业务/个性，绘制的其新的X帖子页面。 [English] Create a new X post page illustrated for Baoyu (refer to the homepage ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@tuzi_ai",
    "tags": [
      "History & Classical Themes",
      "Classical",
      "Social",
      "History"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-234",
    "locale": "zh",
    "title": "朱元璋登基后的推特主页",
    "category": "历史与古典主题",
    "description": "[中文] 创建一个明朝朱元璋登基之后的X帖子页面 [English] Create an X post page of Zhu Yuanzhang after his ascension to the throne in the Ming Dynasty",
    "promptPreview": "[中文] 创建一个明朝朱元璋登基之后的X帖子页面 [English] Create an X post page of Zhu Yuanzhang after his ascension to the throne in the Ming Dynasty",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "History & Classical Themes",
      "Classical",
      "Social",
      "History"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-184",
    "locale": "zh",
    "title": "杜甫朋友圈吐槽茅屋被掀翻",
    "category": "历史与古典主题",
    "description": "[中文] 杜甫发朋友圈吐槽房顶被风刮没了 [English] Du Fu posting on WeChat Moments complaining about his roof being blown away by the wind",
    "promptPreview": "[中文] 杜甫发朋友圈吐槽房顶被风刮没了 [English] Du Fu posting on WeChat Moments complaining about his roof being blown away by the wind",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "History & Classical Themes",
      "History",
      "Tech",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-185",
    "locale": "zh",
    "title": "武则天发微博自拍太魔性了",
    "category": "历史与古典主题",
    "description": "[中文] 武则天自拍登记发微博 [English] Wu Zetian taking a selfie, registering and posting on Weibo.",
    "promptPreview": "[中文] 武则天自拍登记发微博 [English] Wu Zetian taking a selfie, registering and posting on Weibo.",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "History & Classical Themes",
      "History",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-205",
    "locale": "zh",
    "title": "皇宫深处的御用快递驿站",
    "category": "历史与古典主题",
    "description": "[中文] 生成一张古代皇宫 × 快递驿站 [English] Generate an ancient imperial palace × express delivery station",
    "promptPreview": "[中文] 生成一张古代皇宫 × 快递驿站 [English] Generate an ancient imperial palace × express delivery station",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@joshesye",
    "tags": [
      "History & Classical Themes",
      "History",
      "History"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-176",
    "locale": "zh",
    "title": "苏轼被贬首日朋友圈曝光",
    "category": "历史与古典主题",
    "description": "[中文] 苏轼被贬第一天小红书截图 [English] Su Shi's first day of exile Xiaohongshu screenshot",
    "promptPreview": "[中文] 苏轼被贬第一天小红书截图 [English] Su Shi's first day of exile Xiaohongshu screenshot",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "History & Classical Themes",
      "UI",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-379",
    "locale": "en",
    "title": "品牌人格漫画信息图",
    "category": "品牌与标志",
    "description": "Using the uploaded logo, create a highly detailed, comic-style infographic poster: “What This Brand Feels Like” GOAL: Turn the brand into a ",
    "promptPreview": "Using the uploaded logo, create a highly detailed, comic-style infographic poster: “What This Brand Feels Like” GOAL: Turn the brand into a living personality a",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@CallumGrey",
    "tags": [
      "Brand & Logos",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-386",
    "locale": "en",
    "title": "品牌包络产品广告",
    "category": "品牌与标志",
    "description": "The Brand Envelope | GPT Image-2 Prompt #89 This takes any product photo and wraps it in your specific brand world. Different product each t",
    "promptPreview": "The Brand Envelope | GPT Image-2 Prompt #89 This takes any product photo and wraps it in your specific brand world. Different product each time. Same brand, eve",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@SRKDAN",
    "tags": [
      "Brand & Logos",
      "Realistic",
      "Product",
      "Brand",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-160",
    "locale": "en",
    "title": "品牌吉祥物设定图",
    "category": "品牌与标志",
    "description": "Generate a set of icons for {argument name=\"device\" default=\"vintage electronic equipment\"} in {argument name=\"style\" default=\"retro skeuomo",
    "promptPreview": "Generate a set of icons for {argument name=\"device\" default=\"vintage electronic equipment\"} in {argument name=\"style\" default=\"retro skeuomorphic style\"}, inclu",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@TanShilong",
    "tags": [
      "Brand & Logos",
      "UI",
      "Brand",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-459",
    "locale": "zh",
    "title": "品牌奶茶 KV 概念海报",
    "category": "品牌与标志",
    "description": "你是一个品牌视觉识别系统、商业广告创意总监、KV海报设计师和高传播品牌视觉生成系统。 请根据用户输入的【现有品牌名称】，自动识别该品牌最具代表性的品牌Logo形象、品牌名称文字识别特征、主推产品、产品包装、品牌色彩、视觉调性、目标人群和广告传播风格，并生成一张符合该品牌气质的概念",
    "promptPreview": "你是一个品牌视觉识别系统、商业广告创意总监、KV海报设计师和高传播品牌视觉生成系统。 请根据用户输入的【现有品牌名称】，自动识别该品牌最具代表性的品牌Logo形象、品牌名称文字识别特征、主推产品、产品包装、品牌色彩、视觉调性、目标人群和广告传播风格，并生成一张符合该品牌气质的概念 KV 海报。 创作定位： - 基于真实",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "Brand & Logos",
      "Poster",
      "Product",
      "Brand",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-186",
    "locale": "zh",
    "title": "品牌视觉识别图",
    "category": "品牌与标志",
    "description": "[中文] 创建一个包含100种不同奇幻RPG物品的10×10网格，以经典像素艺术风格渲染（16位或32位精灵图美学，让人联想到SNES/GBA时代的日式RPG）。每个物品应出现在其独立的方形瓷砖中，下方带有简短清晰的标签。在白色背景上保持网格整洁。使每个物品在视觉上都有所区分，并",
    "promptPreview": "[中文] 创建一个包含100种不同奇幻RPG物品的10×10网格，以经典像素艺术风格渲染（16位或32位精灵图美学，让人联想到SNES/GBA时代的日式RPG）。每个物品应出现在其独立的方形瓷砖中，下方带有简短清晰的标签。在白色背景上保持网格整洁。使每个物品在视觉上都有所区分，并且每个标签拼写正确。使用清晰的像素边缘、",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@ProperPrompter",
    "tags": [
      "Brand & Logos",
      "Brand",
      "3D",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-136",
    "locale": "en",
    "title": "品牌视觉识别图",
    "category": "品牌与标志",
    "description": "{ \"type\": \"e-commerce landing page hero section\", \"brand\": \"{argument name=\\\"brand name\\\" default=\\\"CLEAR RESET\\\"}\", \"theme\": \"refreshing sk",
    "promptPreview": "{ \"type\": \"e-commerce landing page hero section\", \"brand\": \"{argument name=\\\"brand name\\\" default=\\\"CLEAR RESET\\\"}\", \"theme\": \"refreshing skincare, clean aesthe",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@ryuya__31",
    "tags": [
      "Brand & Logos",
      "Product",
      "Brand",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-342",
    "locale": "en",
    "title": "四季包装 Campaign 宫格",
    "category": "品牌与标志",
    "description": "PHASE 1 - PRODUCT: [ITEM] in [MATERIAL] packaging, minimal label design PHASE 2 - GRID: 2x2 seasonal grid, four distinct brand worlds PHASE ",
    "promptPreview": "PHASE 1 - PRODUCT: [ITEM] in [MATERIAL] packaging, minimal label design PHASE 2 - GRID: 2x2 seasonal grid, four distinct brand worlds PHASE 3 - COMPOSITION: eac",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@SRKDAN",
    "tags": [
      "Brand & Logos",
      "Product",
      "Brand",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-478",
    "locale": "en",
    "title": "夹层式品牌编辑海报",
    "category": "品牌与标志",
    "description": "[BRAND NAME]. You are a world-class editorial designer. STEP 1, DYNAMIC SUBJECT LOGIC: - Subject pick: independently study [BRAND NAME] and ",
    "promptPreview": "[BRAND NAME]. You are a world-class editorial designer. STEP 1, DYNAMIC SUBJECT LOGIC: - Subject pick: independently study [BRAND NAME] and choose the right her",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@iamaiistudio",
    "tags": [
      "Brand & Logos",
      "Poster",
      "Realistic",
      "Product",
      "Commerce",
      "Food",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-138",
    "locale": "en",
    "title": "封面排版设计图",
    "category": "品牌与标志",
    "description": "{ \"type\": \"fashion product catalog layout\", \"theme\": \"A cohesive fashion collection featuring a specific pattern: {argument name=\\\"pattern d",
    "promptPreview": "{ \"type\": \"fashion product catalog layout\", \"theme\": \"A cohesive fashion collection featuring a specific pattern: {argument name=\\\"pattern description\\\" default",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@aiehon_aya",
    "tags": [
      "Brand & Logos",
      "Poster",
      "Product",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-135",
    "locale": "en",
    "title": "应用界面样机图",
    "category": "品牌与标志",
    "description": "{ \"type\": \"website landing page mockup\", \"theme\": \"men's skincare, sleek, professional, dark mode\", \"color_palette\": \"{argument name=\\\"color",
    "promptPreview": "{ \"type\": \"website landing page mockup\", \"theme\": \"men's skincare, sleek, professional, dark mode\", \"color_palette\": \"{argument name=\\\"color scheme\\\" default=\\\"",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@ryuya__31",
    "tags": [
      "Brand & Logos",
      "UI",
      "Product",
      "Brand",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-244",
    "locale": "zh",
    "title": "杜蕾斯茶颜悦色联名海报设计",
    "category": "品牌与标志",
    "description": "[中文] 设计一套杜蕾斯和茶颜悦色联名的宣传物料 [English] Design a set of promotional materials for a Durex and Chayan Yuese co-branding campaign.",
    "promptPreview": "[中文] 设计一套杜蕾斯和茶颜悦色联名的宣传物料 [English] Design a set of promotional materials for a Durex and Chayan Yuese co-branding campaign.",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@akokoi1",
    "tags": [
      "Brand & Logos",
      "Poster",
      "Brand",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-496",
    "locale": "en",
    "title": "水雕品牌 Logo 六宫格",
    "category": "品牌与标志",
    "description": "Create a premium 3x2 grid collage of iconic global brand logos recreated entirely from dynamic water formations, floating above a crystal-cl",
    "promptPreview": "Create a premium 3x2 grid collage of iconic global brand logos recreated entirely from dynamic water formations, floating above a crystal-clear ocean under a vi",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@AIwithSynthia",
    "tags": [
      "Brand & Logos",
      "UI",
      "Realistic",
      "Brand",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-137",
    "locale": "zh",
    "title": "界面交互设计图",
    "category": "品牌与标志",
    "description": "{ \"type\": \"e-commerce landing page hero section mockup\", \"aesthetic\": \"clean, bright, airy, feminine, floral accents with purple flowers, {a",
    "promptPreview": "{ \"type\": \"e-commerce landing page hero section mockup\", \"aesthetic\": \"clean, bright, airy, feminine, floral accents with purple flowers, {argument name=\\\"prima",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@ryuya__31",
    "tags": [
      "Brand & Logos",
      "UI",
      "Product",
      "Brand",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-134",
    "locale": "en",
    "title": "界面交互设计图",
    "category": "品牌与标志",
    "description": "{ \"type\": \"skincare e-commerce landing page mockup\", \"brand\": \"{argument name=\\\"brand name\\\" default=\\\"DERMA CALM\\\"}\", \"color_palette\": [\"wh",
    "promptPreview": "{ \"type\": \"skincare e-commerce landing page mockup\", \"brand\": \"{argument name=\\\"brand name\\\" default=\\\"DERMA CALM\\\"}\", \"color_palette\": [\"white\", \"light blue\", ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@ryuya__31",
    "tags": [
      "Brand & Logos",
      "UI",
      "Product",
      "Brand",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-132",
    "locale": "en",
    "title": "界面交互设计图",
    "category": "品牌与标志",
    "description": "{ \"type\": \"18-panel brand identity and character design document\", \"brand\": { \"name\": \"{argument name=\\\"brand name\\\" default=\\\"沐阳 MUYANG TEA",
    "promptPreview": "{ \"type\": \"18-panel brand identity and character design document\", \"brand\": { \"name\": \"{argument name=\\\"brand name\\\" default=\\\"沐阳 MUYANG TEA\\\"}\", \"industry\": \"{",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Colin_Leeee",
    "tags": [
      "Brand & Logos",
      "UI",
      "Infographic",
      "Realistic",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-130",
    "locale": "en",
    "title": "界面交互设计图",
    "category": "品牌与标志",
    "description": "{ \"type\": \"brand identity and merchandise design board\", \"theme\": { \"color_palette\": \"{argument name=\\\"theme color\\\" default=\\\"pastel pink\\\"",
    "promptPreview": "{ \"type\": \"brand identity and merchandise design board\", \"theme\": { \"color_palette\": \"{argument name=\\\"theme color\\\" default=\\\"pastel pink\\\"} and white\", \"motif",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@chi_vc_",
    "tags": [
      "Brand & Logos",
      "UI",
      "Poster",
      "Illustration",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-294",
    "locale": "zh",
    "title": "精美潮汕菜馆菜单图",
    "category": "品牌与标志",
    "description": "[中文] 生成一张潮菜馆菜单图 [English] Generate a Teochew restaurant menu image.",
    "promptPreview": "[中文] 生成一张潮菜馆菜单图 [English] Generate a Teochew restaurant menu image.",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "Brand & Logos",
      "Brand",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-172",
    "locale": "zh",
    "title": "赛博科幻桃太郎主视觉图",
    "category": "品牌与标志",
    "description": "[中文] 设计虚构动画的钥匙视觉图。主题是「科幻桃太郎」。设计有魅力的角色、背景、标志和宣传语，以一幅美丽插画的形式完成，让世界观在一张图中传达出来。 [English] Design a key visual for a fictional animation. The the",
    "promptPreview": "[中文] 设计虚构动画的钥匙视觉图。主题是「科幻桃太郎」。设计有魅力的角色、背景、标志和宣传语，以一幅美丽插画的形式完成，让世界观在一张图中传达出来。 [English] Design a key visual for a fictional animation. The theme is \"Sci-Fi Momota",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@SSSS_CRYPTOMAN",
    "tags": [
      "Brand & Logos",
      "Illustration",
      "Brand",
      "Character",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-247",
    "locale": "zh",
    "title": "运动健身图标字体设计",
    "category": "品牌与标志",
    "description": "[中文] 生成一套运动类app的iconfont [English] Generate a set of iconfont for a sports app",
    "promptPreview": "[中文] 生成一套运动类app的iconfont [English] Generate a set of iconfont for a sports app",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@akokoi1",
    "tags": [
      "Brand & Logos",
      "Brand",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-444",
    "locale": "zh",
    "title": "迪斯科镜面 3D App 图标",
    "category": "品牌与标志",
    "description": "为【品牌名】生成一个高级 3D App 图标，圆角方形底板，玻璃与金属铬材质，迪斯科球镜面马赛克小方块质感，闪亮高光，柔和工作室灯光，干净极简背景，高端产品图标风格，Blender 3D 渲染，超精细 英文版： A premium 3D app icon for 【Product",
    "promptPreview": "为【品牌名】生成一个高级 3D App 图标，圆角方形底板，玻璃与金属铬材质，迪斯科球镜面马赛克小方块质感，闪亮高光，柔和工作室灯光，干净极简背景，高端产品图标风格，Blender 3D 渲染，超精细 英文版： A premium 3D app icon for 【Product Name】, rounded squa",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@vista8",
    "tags": [
      "Brand & Logos",
      "Product",
      "Brand",
      "3D",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-216",
    "locale": "zh",
    "title": "雅致图案四款时尚单品设计",
    "category": "品牌与标志",
    "description": "[中文] 使用附图中的图案，由专业设计师打造 4 款时尚单品，采用不同的色彩搭配与排版设计，附带穿搭效果图。以雅致的构图凸显图案的美感。格式为 2:3，希望将图像生成模型从 duct-tape-1 指定为 duct-tape-2、3。 [English] Use the patt",
    "promptPreview": "[中文] 使用附图中的图案，由专业设计师打造 4 款时尚单品，采用不同的色彩搭配与排版设计，附带穿搭效果图。以雅致的构图凸显图案的美感。格式为 2:3，希望将图像生成模型从 duct-tape-1 指定为 duct-tape-2、3。 [English] Use the patterns in the attached",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@aiehon_aya",
    "tags": [
      "Brand & Logos",
      "Brand",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-424",
    "locale": "en",
    "title": "FMCG 棒棒糖霓虹广告",
    "category": "商品与电商",
    "description": "Hyper-realistic cinematic FMCG billboard advertising poster for Chupa Chups India, focusing on playful energy, bold flavor explosion, and Ge",
    "promptPreview": "Hyper-realistic cinematic FMCG billboard advertising poster for Chupa Chups India, focusing on playful energy, bold flavor explosion, and Gen-Z candy culture. S",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Diplomeme",
    "tags": [
      "Products & E-commerce",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-441",
    "locale": "en",
    "title": "WILDCAMP 巨型帐篷广告海报",
    "category": "商品与电商",
    "description": "An outdoor adventure advertisement poster featuring a rugged bearded man in full hiking gear standing confidently beside a massive orange ca",
    "promptPreview": "An outdoor adventure advertisement poster featuring a rugged bearded man in full hiking gear standing confidently beside a massive orange camping tent three tim",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Strength04_X",
    "tags": [
      "Products & E-commerce",
      "Poster",
      "Realistic",
      "Brand",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-144",
    "locale": "en",
    "title": "主题海报版式设计",
    "category": "商品与电商",
    "description": "A luxurious cosmetic product advertisement featuring a single elegant glass jar with a shiny gold lid resting on a round, light-colored marb",
    "promptPreview": "A luxurious cosmetic product advertisement featuring a single elegant glass jar with a shiny gold lid resting on a round, light-colored marble slab. The jar has",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@panchaaan_2",
    "tags": [
      "Products & E-commerce",
      "Poster",
      "Product",
      "Brand",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-178",
    "locale": "zh",
    "title": "亚马逊详情图设计",
    "category": "商品与电商",
    "description": "[中文] 生成一套亚马逊 A+=详情图 [English] Generate a set of Amazon A+= detail images",
    "promptPreview": "[中文] 生成一套亚马逊 A+=详情图 [English] Generate a set of Amazon A+= detail images",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@xin_pai88825",
    "tags": [
      "Products & E-commerce",
      "Products",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-155",
    "locale": "en",
    "title": "人物角色设定图",
    "category": "商品与电商",
    "description": "Create {argument name=\"items\" default=\"fan goods\"} for a standard {argument name=\"character type\" default=\"Vtuber\"} in {argument name=\"style",
    "promptPreview": "Create {argument name=\"items\" default=\"fan goods\"} for a standard {argument name=\"character type\" default=\"Vtuber\"} in {argument name=\"style\" default=\"live-acti",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@wtry1102",
    "tags": [
      "Products & E-commerce",
      "Character",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-54",
    "locale": "zh",
    "title": "人物角色设定图",
    "category": "商品与电商",
    "description": "{ \"type\": \"4-panel satirical product advertisement grid\", \"layout\": { \"grid\": \"2x2\", \"panels\": [ { \"position\": \"top-left\", \"product_name\": \"",
    "promptPreview": "{ \"type\": \"4-panel satirical product advertisement grid\", \"layout\": { \"grid\": \"2x2\", \"panels\": [ { \"position\": \"top-left\", \"product_name\": \"{argument name=\\\"top",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@fukumy_ai",
    "tags": [
      "Products & E-commerce",
      "Product",
      "Character",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-475",
    "locale": "en",
    "title": "企鹅造型包装结构板",
    "category": "商品与电商",
    "description": "Using the attached image, create an illustration sheet of professional industrial design packaging for the package (PACKAGE TYPE). A centere",
    "promptPreview": "Using the attached image, create an illustration sheet of professional industrial design packaging for the package (PACKAGE TYPE). A centered heroic 3D renderin",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@iamaiistudio",
    "tags": [
      "Products & E-commerce",
      "Realistic",
      "Illustration",
      "Product",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-194",
    "locale": "zh",
    "title": "健身蛋白粉电商详情页",
    "category": "商品与电商",
    "description": "[中文] 健身蛋白粉电商详情图 [English] Fitness protein powder e-commerce detail image",
    "promptPreview": "[中文] 健身蛋白粉电商详情图 [English] Fitness protein powder e-commerce detail image",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "Products & E-commerce",
      "Products",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-190",
    "locale": "zh",
    "title": "全自动咖啡机产品展示",
    "category": "商品与电商",
    "description": "[中文] 全自动咖啡机电商详情图 [English] Fully automatic coffee machine e-commerce detail image",
    "promptPreview": "[中文] 全自动咖啡机电商详情图 [English] Fully automatic coffee machine e-commerce detail image",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "Products & E-commerce",
      "Products",
      "Tech",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-142",
    "locale": "en",
    "title": "写实摄影风格创作",
    "category": "商品与电商",
    "description": "{ \"type\": \"anime idol merchandise catalog flyer\", \"theme_colors\": \"{argument name=\\\"theme color\\\" default=\\\"pastel blue and pink\\\"}\", \"chara",
    "promptPreview": "{ \"type\": \"anime idol merchandise catalog flyer\", \"theme_colors\": \"{argument name=\\\"theme color\\\" default=\\\"pastel blue and pink\\\"}\", \"character\": { \"name\": \"{a",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@anemone_sd",
    "tags": [
      "Products & E-commerce",
      "UI",
      "Realistic",
      "Brand",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-150",
    "locale": "en",
    "title": "品牌徽标设计图",
    "category": "商品与电商",
    "description": "A bright, summery commercial product photography shot featuring a refreshing beverage on a weathered wooden table. In the sharp foreground, ",
    "promptPreview": "A bright, summery commercial product photography shot featuring a refreshing beverage on a weathered wooden table. In the sharp foreground, there is 1 tall glas",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@highball_cho",
    "tags": [
      "Products & E-commerce",
      "Poster",
      "Realistic",
      "Product",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-462",
    "locale": "zh",
    "title": "复古日系迷你橡皮商品包装",
    "category": "商品与电商",
    "description": "添付されたキャラクターシートをSTRICTなデザインリファレンスとして使用すること。 キャラクターの顔、髪型、目の形、プロポーションは絶対に変更しない。 ■目的： キャラクターを日本の「ちび消しゴム商品」として完全に商品化し、 実際に文房具売り場やガチャで販売されているようなリア",
    "promptPreview": "添付されたキャラクターシートをSTRICTなデザインリファレンスとして使用すること。 キャラクターの顔、髪型、目の形、プロポーションは絶対に変更しない。 ■目的： キャラクターを日本の「ちび消しゴム商品」として完全に商品化し、 実際に文房具売り場やガチャで販売されているようなリアルなパッケージ商品写真を作成する。 ■コ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@ZetoGroovin",
    "tags": [
      "Products & E-commerce",
      "Realistic",
      "Product",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-237",
    "locale": "zh",
    "title": "夏日柑橘苏打高转化广告图",
    "category": "商品与电商",
    "description": "[中文] 图像生成: 商品广告照片, 适合夏天的季节商品, 碳酸饮料, 名称=\"夏柑SODA\", 形状=PET瓶500ml, 研究2025年作为饮料广告的高CTA设计后设计并生成图像规格, 宽高比3:4 [English] Image generation: Product ad",
    "promptPreview": "[中文] 图像生成: 商品广告照片, 适合夏天的季节商品, 碳酸饮料, 名称=\"夏柑SODA\", 形状=PET瓶500ml, 研究2025年作为饮料广告的高CTA设计后设计并生成图像规格, 宽高比3:4 [English] Image generation: Product advertising photo, Sea",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@old_pgmrs_will",
    "tags": [
      "Products & E-commerce",
      "UI",
      "Realistic",
      "Product",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-449",
    "locale": "en",
    "title": "奢华机械腕表技术图鉴",
    "category": "商品与电商",
    "description": "2x2 grid 16:9, do this for 4 most expensive strangest watches ever made: class Haute_Horlogerie_DNA: def __init__(self): self.subject = \"[TI",
    "promptPreview": "2x2 grid 16:9, do this for 4 most expensive strangest watches ever made: class Haute_Horlogerie_DNA: def __init__(self): self.subject = \"[TIMEPIECE]\" self.paren",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Gdgtify",
    "tags": [
      "Products & E-commerce",
      "Infographic",
      "Brand",
      "3D",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-455",
    "locale": "en",
    "title": "巨型舒适洞洞鞋 Campaign",
    "category": "商品与电商",
    "description": "Hyper-realistic premium product advertisement: an oversized futuristic comfort clog sits on a smooth glossy reflective floor. A modern model",
    "promptPreview": "Hyper-realistic premium product advertisement: an oversized futuristic comfort clog sits on a smooth glossy reflective floor. A modern model in soft neutral-ton",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@iamaiistudio",
    "tags": [
      "Products & E-commerce",
      "Poster",
      "Realistic",
      "Product",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-47",
    "locale": "zh",
    "title": "建筑空间场景图",
    "category": "商品与电商",
    "description": "{ \"type\": \"2x2 grid of Japanese digital advertisement banners\", \"layout\": { \"structure\": \"4 equal quadrants\", \"quadrants\": [ { \"position\": \"",
    "promptPreview": "{ \"type\": \"2x2 grid of Japanese digital advertisement banners\", \"layout\": { \"structure\": \"4 equal quadrants\", \"quadrants\": [ { \"position\": \"top-left\", \"theme\": ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@makaneko_AI",
    "tags": [
      "Products & E-commerce",
      "UI",
      "Product",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-454",
    "locale": "en",
    "title": "旅行美食薯片广告海报",
    "category": "商品与电商",
    "description": "Ultra-detailed premium travel-food advertisement poster for [CITY/COUNTRY], vertical composition, inspired by luxury Lay’s-style chips adver",
    "promptPreview": "Ultra-detailed premium travel-food advertisement poster for [CITY/COUNTRY], vertical composition, inspired by luxury Lay’s-style chips advertising. A realistic ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Naiknelofar788",
    "tags": [
      "Products & E-commerce",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-265",
    "locale": "zh",
    "title": "日式潮流广告四联画",
    "category": "商品与电商",
    "description": "[中文] 生成四张虚构的日式广告图片，涵盖不同类型并排排列。采用专业设计师创作的潮流设计。宽高比为1:1 [English] Generate four fictional Japanese advertisement images, covering different typ",
    "promptPreview": "[中文] 生成四张虚构的日式广告图片，涵盖不同类型并排排列。采用专业设计师创作的潮流设计。宽高比为1:1 [English] Generate four fictional Japanese advertisement images, covering different types arranged side by ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@midori_tatsuta",
    "tags": [
      "Products & E-commerce",
      "Poster",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-485",
    "locale": "en",
    "title": "时尚目录电商拼贴",
    "category": "商品与电商",
    "description": "Stylish fashion catalog shoot blending streetwear and luxury branding. Female model wearing burgundy slim-fit top and ivory tailored pants, ",
    "promptPreview": "Stylish fashion catalog shoot blending streetwear and luxury branding. Female model wearing burgundy slim-fit top and ivory tailored pants, posed in confident r",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Mind_Boticni",
    "tags": [
      "Products & E-commerce",
      "UI",
      "Realistic",
      "Brand",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-327",
    "locale": "zh",
    "title": "沉香玫瑰悬浮幻景",
    "category": "商品与电商",
    "description": "[中文] { \"master_prompt_type\": \"超精细8K AI图像生成\", \"global_settings\": { \"resolution\": \"8K UHD\", \"aspect_ratio\": \"2:3 竖版\", \"render_quality\": \"极致锐度、",
    "promptPreview": "[中文] { \"master_prompt_type\": \"超精细8K AI图像生成\", \"global_settings\": { \"resolution\": \"8K UHD\", \"aspect_ratio\": \"2:3 竖版\", \"render_quality\": \"极致锐度、超微细节、电影级光效\", \"style\"",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@meng_dagg695",
    "tags": [
      "Products & E-commerce",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-189",
    "locale": "zh",
    "title": "清新夏日女装连衣裙电商展示",
    "category": "商品与电商",
    "description": "[中文] 夏季女裙电商详情图 [English] Summer women's dress e-commerce detail image",
    "promptPreview": "[中文] 夏季女裙电商详情图 [English] Summer women's dress e-commerce detail image",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "Products & E-commerce",
      "Products",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-318",
    "locale": "zh",
    "title": "珊瑚色极简影棚时尚商业大片",
    "category": "商品与电商",
    "description": "[中文] 超写实高端时尚商业广告大片，使用上传的模特照片作为严格的身份参考。保留精确的面部特征、比例和自然皮肤纹理——无修图，无变形。场景：珊瑚色单色工作室盒，配有光泽反光棋盘格或极简抛光地板。拥有柔和光线渐变的干净几何墙壁。产品：产品放置在前景中心超大位置，因广角透视而占据画面",
    "promptPreview": "[中文] 超写实高端时尚商业广告大片，使用上传的模特照片作为严格的身份参考。保留精确的面部特征、比例和自然皮肤纹理——无修图，无变形。场景：珊瑚色单色工作室盒，配有光泽反光棋盘格或极简抛光地板。拥有柔和光线渐变的干净几何墙壁。产品：产品放置在前景中心超大位置，因广角透视而占据画面主导地位。包装超清晰，文字完全可读，具有",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@Maercihh",
    "tags": [
      "Products & E-commerce",
      "Realistic",
      "Product",
      "Brand",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-438",
    "locale": "en",
    "title": "珠宝微缩城市广告海报",
    "category": "商品与电商",
    "description": "Create a hyper-detailed luxury advertising poster in a cinematic miniature-world style. A gigantic royal diamond necklace with intricate gol",
    "promptPreview": "Create a hyper-detailed luxury advertising poster in a cinematic miniature-world style. A gigantic royal diamond necklace with intricate gold filigree and massi",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Umar__786Ai",
    "tags": [
      "Products & E-commerce",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-192",
    "locale": "zh",
    "title": "电商商品展示图",
    "category": "商品与电商",
    "description": "[中文] AI智能眼镜电商详情图 [English] AI smart glasses e-commerce detail image",
    "promptPreview": "[中文] AI智能眼镜电商详情图 [English] AI smart glasses e-commerce detail image",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "Products & E-commerce",
      "Product",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-313",
    "locale": "zh",
    "title": "电商商品展示设计",
    "category": "商品与电商",
    "description": "[中文] { \"style\": \"超写实奢华化妆品产品摄影\", \"composition\": { \"color_scheme\": \"戏剧性的单色蓝紫色\", \"resolution\": \"8K超高分辨率\", \"depth\": \"电影级景深\", \"aesthetic\": \"高端香氛护",
    "promptPreview": "[中文] { \"style\": \"超写实奢华化妆品产品摄影\", \"composition\": { \"color_scheme\": \"戏剧性的单色蓝紫色\", \"resolution\": \"8K超高分辨率\", \"depth\": \"电影级景深\", \"aesthetic\": \"高端香氛护肤品广告风格\" }, \"product\"",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@Fujimoto_hina",
    "tags": [
      "Products & E-commerce",
      "Poster",
      "Realistic",
      "Product",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-157",
    "locale": "en",
    "title": "电商商品展示设计",
    "category": "商品与电商",
    "description": "{ \"type\": \"e-commerce product infographic\", \"theme\": \"dark mode with {argument name=\\\"accent color\\\" default=\\\"orange\\\"} accents\", \"product\"",
    "promptPreview": "{ \"type\": \"e-commerce product infographic\", \"theme\": \"dark mode with {argument name=\\\"accent color\\\" default=\\\"orange\\\"} accents\", \"product\": { \"brand\": \"{argum",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@AmberPromptai",
    "tags": [
      "Products & E-commerce",
      "UI",
      "Infographic",
      "Product",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-141",
    "locale": "en",
    "title": "电商商品展示设计",
    "category": "商品与电商",
    "description": "{ \"type\": \"promotional banner design set\", \"theme\": \"strawberry advertisement campaign\", \"style\": \"anime illustration, bright, cheerful, com",
    "promptPreview": "{ \"type\": \"promotional banner design set\", \"theme\": \"strawberry advertisement campaign\", \"style\": \"anime illustration, bright, cheerful, commercial graphic desi",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@takadtmnu",
    "tags": [
      "Products & E-commerce",
      "Illustration",
      "Product",
      "Character",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-33",
    "locale": "en",
    "title": "电商商品展示设计",
    "category": "商品与电商",
    "description": "A 3D render of a cute kawaii {argument name=\"subject\" default=\"cloud\"} character on a pure white background. The character has a soft, matte",
    "promptPreview": "A 3D render of a cute kawaii {argument name=\"subject\" default=\"cloud\"} character on a pure white background. The character has a soft, matte, squishy texture re",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@yurunekofree",
    "tags": [
      "Products & E-commerce",
      "UI",
      "Product",
      "Character",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-151",
    "locale": "zh",
    "title": "界面交互设计图",
    "category": "商品与电商",
    "description": "{ \"type\": \"2x2 advertising banner grid\", \"layout\": \"4 distinct quadrants, each featuring a different industry advertisement\", \"quadrants\": [",
    "promptPreview": "{ \"type\": \"2x2 advertising banner grid\", \"layout\": \"4 distinct quadrants, each featuring a different industry advertisement\", \"quadrants\": [ { \"position\": \"top-",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@kitune_fire45",
    "tags": [
      "Products & E-commerce",
      "UI",
      "Brand",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-301",
    "locale": "zh",
    "title": "终结者机器人淘宝详情页",
    "category": "商品与电商",
    "description": "[中文] 生成图片: T-800机器人的淘宝商品详情页，展示: 机器人的正面侧面背面三视图， 产品价格， 产品细节， 功能和使用场景等 [English] Generate image: Taobao product detail page of a T-800 robot, s",
    "promptPreview": "[中文] 生成图片: T-800机器人的淘宝商品详情页，展示: 机器人的正面侧面背面三视图， 产品价格， 产品细节， 功能和使用场景等 [English] Generate image: Taobao product detail page of a T-800 robot, showing: front, side,",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@rionaifantasy",
    "tags": [
      "Products & E-commerce",
      "Product",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-264",
    "locale": "zh",
    "title": "美妆产品广告图",
    "category": "商品与电商",
    "description": "[中文] 为Z世代设计的可爱Y2K风格的平价化妆品广告图像。使用鲜艳的配色，包括荧光色。纵横比为3:4。 [English] Cute Y2K style affordable cosmetics advertising image designed for Gen Z. Usi",
    "promptPreview": "[中文] 为Z世代设计的可爱Y2K风格的平价化妆品广告图像。使用鲜艳的配色，包括荧光色。纵横比为3:4。 [English] Cute Y2K style affordable cosmetics advertising image designed for Gen Z. Using vibrant color sch",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@midori_tatsuta",
    "tags": [
      "Products & E-commerce",
      "Products",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-358",
    "locale": "en",
    "title": "草莓能量饮料商业广告",
    "category": "商品与电商",
    "description": "A hyper-realistic commercial advertisement blending energy drink and sports branding. A dynamic athletic woman mid-air jump, wearing modern ",
    "promptPreview": "A hyper-realistic commercial advertisement blending energy drink and sports branding. A dynamic athletic woman mid-air jump, wearing modern sportswear (light tr",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@SPEEDAI07",
    "tags": [
      "Products & E-commerce",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-317",
    "locale": "zh",
    "title": "震撼视觉的深红影棚广角美妆大片",
    "category": "商品与电商",
    "description": "[中文] 照片级真实感的大胆美妆宣传活动，使用上传的模特作为精确的身份参考。不做面部改变，不做平滑处理。 场景：深红色饱和的摄影棚环境，具有高对比度的地板图案或光滑表面。 产品：产品被握持或放置在极其靠近镜头的位置，由于透视关系显得巨大。 模特姿势：俏皮或自信的微笑，手臂完全伸向",
    "promptPreview": "[中文] 照片级真实感的大胆美妆宣传活动，使用上传的模特作为精确的身份参考。不做面部改变，不做平滑处理。 场景：深红色饱和的摄影棚环境，具有高对比度的地板图案或光滑表面。 产品：产品被握持或放置在极其靠近镜头的位置，由于透视关系显得巨大。 模特姿势：俏皮或自信的微笑，手臂完全伸向相机，手指因广角镜头而略微变形。透过太阳",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@Maercihh",
    "tags": [
      "Products & E-commerce",
      "Realistic",
      "Product",
      "Brand",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-333",
    "locale": "zh",
    "title": "AI 眼镜爆炸拆解图",
    "category": "图表与信息可视化",
    "description": "生成一张AI眼镜的爆炸视图，包含每个组件的名称以及这款产品的几大核心卖点。",
    "promptPreview": "生成一张AI眼镜的爆炸视图，包含每个组件的名称以及这款产品的几大核心卖点。",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "苍何原创实测（公众号文章《我逆向了 329 条 GPT-Image2 提示词模板，全部开源！》）",
    "tags": [
      "Charts & Infographics",
      "Charts",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-407",
    "locale": "en",
    "title": "Neuro-AI 混合系统信息图",
    "category": "图表与信息可视化",
    "description": "Create a premium square “neuro-AI hybrid system infographic” designed as a scientific cognitive engineering handbook page. Visual Direction:",
    "promptPreview": "Create a premium square “neuro-AI hybrid system infographic” designed as a scientific cognitive engineering handbook page. Visual Direction: • 1:1 composition •",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@YaZoraiz",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-183",
    "locale": "zh",
    "title": "一张中文健身信息图",
    "category": "图表与信息可视化",
    "description": "请生成一张中文健身信息图，主题为：【xxx】。 要求这张图既专业又实用，适合普通成年人作为训练参考。默认对象为无严重伤病的健康成年人；如果没有额外说明，默认训练目标为“增肌 + 基础力量提升”，默认训练水平为“新手到中级之间”，默认训练场景为“普通健身房”，默认单次训练时长控制在",
    "promptPreview": "请生成一张中文健身信息图，主题为：【xxx】。 要求这张图既专业又实用，适合普通成年人作为训练参考。默认对象为无严重伤病的健康成年人；如果没有额外说明，默认训练目标为“增肌 + 基础力量提升”，默认训练水平为“新手到中级之间”，默认训练场景为“普通健身房”，默认单次训练时长控制在 40–60 分钟内。 请根据【训练主题",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "Charts & Infographics",
      "Infographic",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-171",
    "locale": "zh",
    "title": "信息图可视化设计",
    "category": "图表与信息可视化",
    "description": "[中文] 创建一个包含 10x10 网格的图像，每个对象名称都以字母 a 开头。 [English] create an image with 10x10 grid of objects that have the names starting with letter a.",
    "promptPreview": "[中文] 创建一个包含 10x10 网格的图像，每个对象名称都以字母 a 开头。 [English] create an image with 10x10 grid of objects that have the names starting with letter a.",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@umesh_ai",
    "tags": [
      "Charts & Infographics",
      "Infographic",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-102",
    "locale": "en",
    "title": "信息图可视化设计",
    "category": "图表与信息可视化",
    "description": "Search the web for {argument name=\"performance description\" default=\"this week’s standout individual performance in Champion’s League\"}, usi",
    "promptPreview": "Search the web for {argument name=\"performance description\" default=\"this week’s standout individual performance in Champion’s League\"}, using exact stats and g",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@maxescu",
    "tags": [
      "Charts & Infographics",
      "Infographic",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-73",
    "locale": "zh",
    "title": "信息图可视化设计",
    "category": "图表与信息可视化",
    "description": "{ \"type\": \"complex urban systems atlas infographic\", \"style\": \"{argument name=\\\"color palette\\\" default=\\\"dark background with glowing blue,",
    "promptPreview": "{ \"type\": \"complex urban systems atlas infographic\", \"style\": \"{argument name=\\\"color palette\\\" default=\\\"dark background with glowing blue, gold, and purple ac",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@hx831126",
    "tags": [
      "Charts & Infographics",
      "Infographic",
      "Illustration",
      "3D",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-72",
    "locale": "en",
    "title": "信息图可视化设计",
    "category": "图表与信息可视化",
    "description": "{ \"type\": \"scientific botanical infographic poster\", \"subject\": \"{argument name=\\\"plant species\\\" default=\\\"Pomegranate (Punica granatum)\\\"}",
    "promptPreview": "{ \"type\": \"scientific botanical infographic poster\", \"subject\": \"{argument name=\\\"plant species\\\" default=\\\"Pomegranate (Punica granatum)\\\"}\", \"style\": \"vintage",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@hx831126",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-70",
    "locale": "en",
    "title": "信息图可视化设计",
    "category": "图表与信息可视化",
    "description": "{ \"type\": \"technical infographic\", \"subject\": \"{argument name=\\\"subject matter\\\" default=\\\"digital photography process\\\"}\", \"header\": { \"tit",
    "promptPreview": "{ \"type\": \"technical infographic\", \"subject\": \"{argument name=\\\"subject matter\\\" default=\\\"digital photography process\\\"}\", \"header\": { \"title\": \"{argument name",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@hx831126",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-69",
    "locale": "en",
    "title": "信息图可视化设计",
    "category": "图表与信息可视化",
    "description": "{ \"type\": \"comprehensive medical infographic\", \"style\": \"highly detailed 3D medical illustration, clinical white background, clean typograph",
    "promptPreview": "{ \"type\": \"comprehensive medical infographic\", \"style\": \"highly detailed 3D medical illustration, clinical white background, clean typography\", \"header\": { \"tit",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@hx831126",
    "tags": [
      "Charts & Infographics",
      "Infographic",
      "Poster",
      "Illustration",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-67",
    "locale": "zh",
    "title": "信息图可视化设计",
    "category": "图表与信息可视化",
    "description": "{ \"type\": \"medical infographic poster\", \"style\": \"highly detailed anatomical illustrations, clean structured layout, scientific diagrammatic",
    "promptPreview": "{ \"type\": \"medical infographic poster\", \"style\": \"highly detailed anatomical illustrations, clean structured layout, scientific diagrammatic style\", \"color_pale",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@hx831126",
    "tags": [
      "Charts & Infographics",
      "Infographic",
      "Poster",
      "Illustration",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-66",
    "locale": "en",
    "title": "信息图可视化设计",
    "category": "图表与信息可视化",
    "description": "{ \"type\": \"fashion design process infographic\", \"title\": \"{argument name=\\\"main title\\\" default=\\\"一件女装诞生的因果链 THE CAUSAL CHAIN OF A WOMEN'S G",
    "promptPreview": "{ \"type\": \"fashion design process infographic\", \"title\": \"{argument name=\\\"main title\\\" default=\\\"一件女装诞生的因果链 THE CAUSAL CHAIN OF A WOMEN'S GARMENT\\\"}\", \"subtitl",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@hx831126",
    "tags": [
      "Charts & Infographics",
      "Infographic",
      "Illustration",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-65",
    "locale": "en",
    "title": "信息图可视化设计",
    "category": "图表与信息可视化",
    "description": "A breathtaking and extremely complex world-building infographic masterpiece conceptualizing the \"{argument name=\"theme\" default=\"Fundamental",
    "promptPreview": "A breathtaking and extremely complex world-building infographic masterpiece conceptualizing the \"{argument name=\"theme\" default=\"Fundamental Differences between",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@GeekCatX",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Illustration",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-64",
    "locale": "zh",
    "title": "信息图可视化设计",
    "category": "图表与信息可视化",
    "description": "{\"type\":\"infographic poster\",\"style\":\"cute flat vector illustration, cozy, warm, soft shading, {argument name=\\\"color palette\\\" default=\\\"pa",
    "promptPreview": "{\"type\":\"infographic poster\",\"style\":\"cute flat vector illustration, cozy, warm, soft shading, {argument name=\\\"color palette\\\" default=\\\"pastel Morandi colors,",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@j_zou93",
    "tags": [
      "Charts & Infographics",
      "Infographic",
      "Poster",
      "Illustration",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-55",
    "locale": "en",
    "title": "信息图可视化设计",
    "category": "图表与信息可视化",
    "description": "Help me create a detailed production flowchart for the dish {argument name=\"dish name\" default=\"Fried Pork with Chili\"}, in a realistic styl",
    "promptPreview": "Help me create a detailed production flowchart for the dish {argument name=\"dish name\" default=\"Fried Pork with Chili\"}, in a realistic style, suitable for Xiao",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Kurt_Rousey466",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Realistic",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-51",
    "locale": "zh",
    "title": "信息图可视化设计",
    "category": "图表与信息可视化",
    "description": "{ \"type\": \"7-day fashion lookbook infographic\", \"header\": { \"title\": \"{argument name=\\\"main title\\\" default=\\\"一周穿搭指南\\\"}\", \"subtitle\": \"{argu",
    "promptPreview": "{ \"type\": \"7-day fashion lookbook infographic\", \"header\": { \"title\": \"{argument name=\\\"main title\\\" default=\\\"一周穿搭指南\\\"}\", \"subtitle\": \"{argument name=\\\"style ke",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@yyyole",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-14",
    "locale": "zh",
    "title": "信息图可视化设计",
    "category": "图表与信息可视化",
    "description": "视觉设计规格描述：画幅比 9:16（竖版手机信息图）；背景纹理为具有呼吸感的米色手工纸（Handmade Washi Paper），带微小纤维纹理，边角有轻微水渍晕染；配色方案为熟番茄红（#E23A28）、初榨橄榄油金黄（#F2C94C）、嫩草绿（#6FCF97）、碳黑墨线；排版",
    "promptPreview": "视觉设计规格描述：画幅比 9:16（竖版手机信息图）；背景纹理为具有呼吸感的米色手工纸（Handmade Washi Paper），带微小纤维纹理，边角有轻微水渍晕染；配色方案为熟番茄红（#E23A28）、初榨橄榄油金黄（#F2C94C）、嫩草绿（#6FCF97）、碳黑墨线；排版逻辑为顶端大标题、中间 Z 字形流线、底",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "小红书号Roy_Jay",
    "tags": [
      "Charts & Infographics",
      "Infographic",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-87",
    "locale": "en",
    "title": "关系图谱信息图",
    "category": "图表与信息可视化",
    "description": "Please generate a high-design character relationship map poster based on {argument name=\"theme\" default=\"Demon Slayer\"}. This image should n",
    "promptPreview": "Please generate a high-design character relationship map poster based on {argument name=\"theme\" default=\"Demon Slayer\"}. This image should not be a simple illus",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@MrLarus",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Poster",
      "Commerce",
      "Social",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-77",
    "locale": "en",
    "title": "关系图谱信息图",
    "category": "图表与信息可视化",
    "description": "Please generate a high-quality vertical \"Popular Science Encyclopedia Image\" based on {argument name=\"theme\" default=\"animals\"}. This image ",
    "promptPreview": "Please generate a high-quality vertical \"Popular Science Encyclopedia Image\" based on {argument name=\"theme\" default=\"animals\"}. This image is not a regular pos",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@alanlovelq",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-71",
    "locale": "zh",
    "title": "关系图谱信息图",
    "category": "图表与信息可视化",
    "description": "{ \"type\": \"technical infographic and exploded view diagram\", \"header\": { \"title\": \"{argument name=\\\"main title\\\" default=\\\"佳能 EOS R5 成像系统剖面 ",
    "promptPreview": "{ \"type\": \"technical infographic and exploded view diagram\", \"header\": { \"title\": \"{argument name=\\\"main title\\\" default=\\\"佳能 EOS R5 成像系统剖面 CANON EOS R5 IMAGING",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@hx831126",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-380",
    "locale": "en",
    "title": "冠状病毒尺度缩放科学信息图",
    "category": "图表与信息可视化",
    "description": "instructions> [SUBJECT]=Coronavirus. A hyper-realistic 3D zoom-sequence infographic generated from a single input: [SUBJECT]. The system aut",
    "promptPreview": "instructions> [SUBJECT]=Coronavirus. A hyper-realistic 3D zoom-sequence infographic generated from a single input: [SUBJECT]. The system auto-detects scale laye",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Gdgtify",
    "tags": [
      "Charts & Infographics",
      "Infographic",
      "Realistic",
      "3D",
      "Tech",
      "Story"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-296",
    "locale": "zh",
    "title": "博物馆级中文拆解信息图鉴",
    "category": "图表与信息可视化",
    "description": "[中文] 请根据【主题】自动生成一张“博物馆图鉴式中文拆解信息图”。 要求整张图兼具真实写实主视觉、结构拆解、中文标注、材质说明、纹样寓意、色彩含义和核心特征总结。你需要根据【主题】自动判断最合适的主体对象、服饰体系、器物结构、时代风格、关键部件、材质工艺、颜色方案与版式结构，用",
    "promptPreview": "[中文] 请根据【主题】自动生成一张“博物馆图鉴式中文拆解信息图”。 要求整张图兼具真实写实主视觉、结构拆解、中文标注、材质说明、纹样寓意、色彩含义和核心特征总结。你需要根据【主题】自动判断最合适的主体对象、服饰体系、器物结构、时代风格、关键部件、材质工艺、颜色方案与版式结构，用户无需再提供其他信息。 整体风格应为：国",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-456",
    "locale": "en",
    "title": "历史事件 2x2 可视化地图",
    "category": "图表与信息可视化",
    "description": "2x2 grid, 16:9, do this for 4 famous historical events: Render_Target = ( 3D_Diorama_Map_Of_[HISTORICAL_EVENT]_In_[CITY] * 1.2 ) + ( First-P",
    "promptPreview": "2x2 grid, 16:9, do this for 4 famous historical events: Render_Target = ( 3D_Diorama_Map_Of_[HISTORICAL_EVENT]_In_[CITY] * 1.2 ) + ( First-Person_Account_Infogr",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Gdgtify",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Classical",
      "Commerce",
      "Travel",
      "History"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-443",
    "locale": "en",
    "title": "塔可爆炸拆解信息图",
    "category": "图表与信息可视化",
    "description": "Create a hyper-realistic exploded vertical infographic composition of tacos. Top → Bottom structure: Fresh Lettuce (crisp green texture with",
    "promptPreview": "Create a hyper-realistic exploded vertical infographic composition of tacos. Top → Bottom structure: Fresh Lettuce (crisp green texture with natural folds) → To",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Strength04_X",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Realistic",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-364",
    "locale": "en",
    "title": "奢华个人色彩档案信息图",
    "category": "图表与信息可视化",
    "description": "LUXURY PERSONAL COLOR PROFILE — EDITORIAL LAYOUT Studio portrait of subject as anchor — skin retouched to luminous glass-like perfection, pr",
    "promptPreview": "LUXURY PERSONAL COLOR PROFILE — EDITORIAL LAYOUT Studio portrait of subject as anchor — skin retouched to luminous glass-like perfection, preserved natural stru",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@meng_dagg695",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Realistic",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-469",
    "locale": "zh",
    "title": "导览式科普绘本",
    "category": "图表与信息可视化",
    "description": "《导览式科普绘本》提示词： 请根据【主题】创作一张高完成度的「导览式科普绘本」风格插画。 这是一张结合“大型场景主视觉 + 导览路线 + 可爱导览 IP + 知识站点 + 儿童科普绘本质感”的场景导览式科普图解页。画面需要让观者像被带着参观一个复杂系统一样，边看边理解主题背后的运",
    "promptPreview": "《导览式科普绘本》提示词： 请根据【主题】创作一张高完成度的「导览式科普绘本」风格插画。 这是一张结合“大型场景主视觉 + 导览路线 + 可爱导览 IP + 知识站点 + 儿童科普绘本质感”的场景导览式科普图解页。画面需要让观者像被带着参观一个复杂系统一样，边看边理解主题背后的运行逻辑、空间结构、流程关系和关键知识点。",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "Charts & Infographics",
      "Infographic",
      "Illustration",
      "Character",
      "Education",
      "Travel",
      "Story"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-248",
    "locale": "zh",
    "title": "景德镇青花瓷全景解说图谱",
    "category": "图表与信息可视化",
    "description": "[中文] 为我生成景德镇青花瓷的详细解说图，配上详细的中文知识解析 [English] Generate a detailed explanatory diagram of Jingdezhen blue and white porcelain, accompanied by d",
    "promptPreview": "[中文] 为我生成景德镇青花瓷的详细解说图，配上详细的中文知识解析 [English] Generate a detailed explanatory diagram of Jingdezhen blue and white porcelain, accompanied by detailed Chinese know",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@joshesye",
    "tags": [
      "Charts & Infographics",
      "Infographic",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-179",
    "locale": "zh",
    "title": "蒸汽朋克射手座解剖图谱",
    "category": "图表与信息可视化",
    "description": "[中文] （Steampunk Scientific Illustrator）你是一位专业复古蒸汽朋克解剖图谱设计师，擅长星座机械结构科普海报。根据用户指定的【{constellation_name}】，生成一张复古蒸汽朋克风格星座解剖图谱海报：顶部标题栏为“{constella",
    "promptPreview": "[中文] （Steampunk Scientific Illustrator）你是一位专业复古蒸汽朋克解剖图谱设计师，擅长星座机械结构科普海报。根据用户指定的【{constellation_name}】，生成一张复古蒸汽朋克风格星座解剖图谱海报：顶部标题栏为“{constellation_name}解剖图谱”或“ANA",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@GeekCatX",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-447",
    "locale": "en",
    "title": "现代地铁工程信息图",
    "category": "图表与信息可视化",
    "description": "Create a premium square “reference-style urban transportation infographic” centered around a futuristic modern metro system called the {METR",
    "promptPreview": "Create a premium square “reference-style urban transportation infographic” centered around a futuristic modern metro system called the {METRO_NAME}, designed as",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@j_smeaton99",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-494",
    "locale": "en",
    "title": "电动巴士工程信息图",
    "category": "图表与信息可视化",
    "description": "Create a premium square “reference-style sustainable transportation infographic” centered around a futuristic electric city bus called the {",
    "promptPreview": "Create a premium square “reference-style sustainable transportation infographic” centered around a futuristic electric city bus called the {E_BUS_NAME}, designe",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@j_smeaton99",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-8",
    "locale": "zh",
    "title": "科普百科图",
    "category": "图表与信息可视化",
    "description": "根据【主题】生成一张高质量竖版「科普百科图」。 这张图不是普通海报，也不是单纯插画，而是一张兼具图鉴感、百科感、信息结构感和收藏感的模块化科普信息图。整体风格参考高级博物图鉴、现代百科书页、生活方式知识卡，以及社交媒体上更容易传播的信息图风格。 让画面包含： 一个清晰好看的主题主",
    "promptPreview": "根据【主题】生成一张高质量竖版「科普百科图」。 这张图不是普通海报，也不是单纯插画，而是一张兼具图鉴感、百科感、信息结构感和收藏感的模块化科普信息图。整体风格参考高级博物图鉴、现代百科书页、生活方式知识卡，以及社交媒体上更容易传播的信息图风格。 让画面包含： 一个清晰好看的主题主视觉 若干局部特征放大细节 多个圆角模块",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "小红书号1055699679",
    "tags": [
      "Charts & Infographics",
      "Infographic",
      "Poster",
      "Illustration",
      "Commerce",
      "Education",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-222",
    "locale": "zh",
    "title": "精致模块化科普百科图鉴",
    "category": "图表与信息可视化",
    "description": "[中文] 请根据【主题】生成一张高质量竖版「科普百科图」。 这张图不是普通海报，也不是单纯插画，而是一张兼具“图鉴感、百科感、信息结构感、收藏感”的模块化科普信息图。整体风格参考高级博物图鉴、现代百科书页、生活方式知识卡和社交媒体高传播信息图的结合。 请让画面包含： - 一个清晰",
    "promptPreview": "[中文] 请根据【主题】生成一张高质量竖版「科普百科图」。 这张图不是普通海报，也不是单纯插画，而是一张兼具“图鉴感、百科感、信息结构感、收藏感”的模块化科普信息图。整体风格参考高级博物图鉴、现代百科书页、生活方式知识卡和社交媒体高传播信息图的结合。 请让画面包含： - 一个清晰漂亮的主题主视觉 - 若干局部特征放大细",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-218",
    "locale": "zh",
    "title": "绘制科学百科知识图谱",
    "category": "图表与信息可视化",
    "description": "[中文] 角色：世界级科学百科插画师兼知识图谱架构师 任务：以经典、无品牌标识（无任何 Logo）的科学百科风格，创作一幅细节极致丰富、结构极其精巧、视觉效果惊艳的「环球图解百科科学信息图」。 题材选择：从【人物、植物、动物】中任选其一。 具体对象：【例如：大王乌贼 / 列奥纳多",
    "promptPreview": "[中文] 角色：世界级科学百科插画师兼知识图谱架构师 任务：以经典、无品牌标识（无任何 Logo）的科学百科风格，创作一幅细节极致丰富、结构极其精巧、视觉效果惊艳的「环球图解百科科学信息图」。 题材选择：从【人物、植物、动物】中任选其一。 具体对象：【例如：大王乌贼 / 列奥纳多・达・芬奇 / 红杉树】 风格：采用复古",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@GeekCatX",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Realistic",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-214",
    "locale": "en",
    "title": "绘制金瓶梅知识图谱",
    "category": "图表与信息可视化",
    "description": "Role: World-class Scientific Encyclopedia Illustrator & Knowledge Graph Architect. Task: Generate a highly detailed, extremely intricate, an",
    "promptPreview": "Role: World-class Scientific Encyclopedia Illustrator & Knowledge Graph Architect. Task: Generate a highly detailed, extremely intricate, and visually stunning ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@xiaoxiaodong01",
    "tags": [
      "Charts & Infographics",
      "UI",
      "Infographic",
      "Realistic",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-210",
    "locale": "zh",
    "title": "萌系大模型训练图解",
    "category": "图表与信息可视化",
    "description": "[中文] 可爱地解释一下大语言模型训练过程 [English] Cute explanation of the large language model training process",
    "promptPreview": "[中文] 可爱地解释一下大语言模型训练过程 [English] Cute explanation of the large language model training process",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@op7418",
    "tags": [
      "Charts & Infographics",
      "Infographic",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-457",
    "locale": "en",
    "title": "运动轨迹舞者光绘海报",
    "category": "图表与信息可视化",
    "description": "<instructions> input = 3 legendary dancer/choreographer input = ai inferred signature movement sequence function render_kinesphere ($ dancer",
    "promptPreview": "<instructions> input = 3 legendary dancer/choreographer input = ai inferred signature movement sequence function render_kinesphere ($ dancer, $ movement) anchor",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Gdgtify",
    "tags": [
      "Charts & Infographics",
      "Poster",
      "Realistic",
      "Character",
      "Tech",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-448",
    "locale": "en",
    "title": "1942 空战街机电影城",
    "category": "场景与叙事",
    "description": "Ultra-realistic cinematic “1942 City” scene, retro WWII arcade-shooter-inspired metropolis in real life, intense aerial-war atmosphere, vert",
    "promptPreview": "Ultra-realistic cinematic “1942 City” scene, retro WWII arcade-shooter-inspired metropolis in real life, intense aerial-war atmosphere, vertical 4:5 editorial c",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Xaroon_x",
    "tags": [
      "Scenes & Storytelling",
      "UI",
      "Realistic",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-477",
    "locale": "en",
    "title": "Instagram 餐桌俯拍创意",
    "category": "场景与叙事",
    "description": "Top-down ultra-wide hyper-realistic shot of 4 real people seated at a square dining table. Pull the camera way back so a wide ring of empty ",
    "promptPreview": "Top-down ultra-wide hyper-realistic shot of 4 real people seated at a square dining table. Pull the camera way back so a wide ring of empty floor surrounds the ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@iamaiistudio",
    "tags": [
      "Scenes & Storytelling",
      "UI",
      "Realistic",
      "Brand",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-493",
    "locale": "en",
    "title": "东京旅行 13 格视频封面",
    "category": "场景与叙事",
    "description": "Generate image of young female travel vlogger exploring Tokyo across multiple candid moments, extremely beautiful with long dark wavy hair, ",
    "promptPreview": "Generate image of young female travel vlogger exploring Tokyo across multiple candid moments, extremely beautiful with long dark wavy hair, wearing stylish Japa",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@AIwithSynthia",
    "tags": [
      "Scenes & Storytelling",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-394",
    "locale": "en",
    "title": "中世纪村庄双精灵冒险者",
    "category": "场景与叙事",
    "description": "Create a cinematic dark-fantasy medieval street scene in ultra-realistic 3D game concept art style, widescreen 16:9. In the foreground, show",
    "promptPreview": "Create a cinematic dark-fantasy medieval street scene in ultra-realistic 3D game concept art style, widescreen 16:9. In the foreground, show two adult elven adv",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@RamonVi25791296",
    "tags": [
      "Scenes & Storytelling",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-422",
    "locale": "en",
    "title": "冬季生存惊悚 Storyboard",
    "category": "场景与叙事",
    "description": "Cinematic Survival Thriller Storyboard Prompt Create a premium cinematic storyboard presentation sheet for a prestige winter survival thrill",
    "promptPreview": "Cinematic Survival Thriller Storyboard Prompt Create a premium cinematic storyboard presentation sheet for a prestige winter survival thriller. Ultra-detailed p",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@zulkarnaimx",
    "tags": [
      "Scenes & Storytelling",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-182",
    "locale": "zh",
    "title": "千禧年日系校园喜剧场景",
    "category": "场景与叙事",
    "description": "[中文] 2000 年代面向中学生的日剧喜剧场景 [English] 2000s Japanese TV drama comedy scene aimed at middle school students",
    "promptPreview": "[中文] 2000 年代面向中学生的日剧喜剧场景 [English] 2000s Japanese TV drama comedy scene aimed at middle school students",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@UminekoStudio",
    "tags": [
      "Scenes & Storytelling",
      "Scenes",
      "Tech",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-419",
    "locale": "en",
    "title": "可颂烘焙流程 Storyboard",
    "category": "场景与叙事",
    "description": "Create a crisp, clean infographic storyboard poster for THE CROISSANT BAKER. Wide 16:9 layout, white background, black borders, bold black t",
    "promptPreview": "Create a crisp, clean infographic storyboard poster for THE CROISSANT BAKER. Wide 16:9 layout, white background, black borders, bold black typography, premium P",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@TechieBySA",
    "tags": [
      "Scenes & Storytelling",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-324",
    "locale": "zh",
    "title": "复古巴士上的红风衣女郎",
    "category": "场景与叙事",
    "description": "[中文] 一位时尚年轻女子坐在老式复古巴士的前缘，身穿红色长风衣、羊毛无檐小便帽、圆形蓝色反光太阳镜、叠层项链和粗犷的棕色皮靴。她有着波浪状金发，带着自信而梦幻的表情，仰望天空。巴士漆面剥落，呈青绿色与铁锈红色调。明亮清澈的蓝天，城市背景建筑极少，柔和日光，电影级色彩分级，浅景深",
    "promptPreview": "[中文] 一位时尚年轻女子坐在老式复古巴士的前缘，身穿红色长风衣、羊毛无檐小便帽、圆形蓝色反光太阳镜、叠层项链和粗犷的棕色皮靴。她有着波浪状金发，带着自信而梦幻的表情，仰望天空。巴士漆面剥落，呈青绿色与铁锈红色调。明亮清澈的蓝天，城市背景建筑极少，柔和日光，电影级色彩分级，浅景深，高端时尚旅行氛围，编辑摄影，超写实，4",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@iamsofiaijaz",
    "tags": [
      "Scenes & Storytelling",
      "UI",
      "Realistic",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-392",
    "locale": "en",
    "title": "头发里的微型城市",
    "category": "场景与叙事",
    "description": "Macro photograph of a miniature city hidden in human hair, clearly on a real human head, with part of the forehead and hairline visible, rea",
    "promptPreview": "Macro photograph of a miniature city hidden in human hair, clearly on a real human head, with part of the forehead and hairline visible, realistic skin texture ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@krafterlab",
    "tags": [
      "Scenes & Storytelling",
      "UI",
      "Realistic",
      "Illustration",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-413",
    "locale": "en",
    "title": "当代舞现场 Storyboard",
    "category": "场景与叙事",
    "description": "Create a raw contemporary dance performance storyboard focused on intense physical movement and live singing. Use reference image for the ch",
    "promptPreview": "Create a raw contemporary dance performance storyboard focused on intense physical movement and live singing. Use reference image for the character. 16:9 storyb",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@ogbenniasamuel2",
    "tags": [
      "Scenes & Storytelling",
      "UI",
      "Realistic",
      "Character",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-340",
    "locale": "zh",
    "title": "彼岸花丛中的红妆女子",
    "category": "场景与叙事",
    "description": "异质感oc，绝美红妆女子，位于彼岸花丛中，张力。 唐琬《钗头凤·世情薄》 世情薄，人情恶，雨送黄昏花易落。晓风干，泪痕残。欲笺心事，独语斜阑。难，难，难！",
    "promptPreview": "异质感oc，绝美红妆女子，位于彼岸花丛中，张力。 唐琬《钗头凤·世情薄》 世情薄，人情恶，雨送黄昏花易落。晓风干，泪痕残。欲笺心事，独语斜阑。难，难，难！",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@xiaofenggan",
    "tags": [
      "Scenes & Storytelling",
      "Scenes",
      "History"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-238",
    "locale": "zh",
    "title": "星云巨鲤与小人的奇幻对话",
    "category": "场景与叙事",
    "description": "[中文] 一幅超现实主义数字插画风格，采用低角度仰拍视角。画面描绘了一条巨型彩色锦鲤遨游在梦幻般的星云中，四周环绕着色彩鲜艳的星云与气泡。 画面中央还站着一个小人，背对观众，神情平静地仰望空中这条巨大的锦鲤，锦鲤头向下看着小人。 整体画面呈现出强烈的大小对比，氛围空灵又梦幻。比例",
    "promptPreview": "[中文] 一幅超现实主义数字插画风格，采用低角度仰拍视角。画面描绘了一条巨型彩色锦鲤遨游在梦幻般的星云中，四周环绕着色彩鲜艳的星云与气泡。 画面中央还站着一个小人，背对观众，神情平静地仰望空中这条巨大的锦鲤，锦鲤头向下看着小人。 整体画面呈现出强烈的大小对比，氛围空灵又梦幻。比例9:16 [English] A sur",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "Scenes & Storytelling",
      "Illustration",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-3"
  },
  {
    "id": "freestylefly-315",
    "locale": "zh",
    "title": "棘龙巨口中的酷飒少女与史前奇观",
    "category": "场景与叙事",
    "description": "[中文] 超写实电影级奇幻场景，设定在郁郁葱葱的史前丛林山谷中。一只巨大的棘龙站在浅河边，它那长而类似鳄鱼的巨颚张得很大。一位年轻女子平静地坐在恐龙张开的嘴里，完美居中，双腿微微向前悬挂。她有一头深色直发，表情镇定无畏，皮肤纹理逼真。她身穿合身的黑色长袖短款上衣，蓝色牛仔短裤和黑",
    "promptPreview": "[中文] 超写实电影级奇幻场景，设定在郁郁葱葱的史前丛林山谷中。一只巨大的棘龙站在浅河边，它那长而类似鳄鱼的巨颚张得很大。一位年轻女子平静地坐在恐龙张开的嘴里，完美居中，双腿微微向前悬挂。她有一头深色直发，表情镇定无畏，皮肤纹理逼真。她身穿合身的黑色长袖短款上衣，蓝色牛仔短裤和黑色及膝战术靴。衣服和腿上可见微小的血迹和",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrDasOnX",
    "tags": [
      "Scenes & Storytelling",
      "Poster",
      "Realistic",
      "Character",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-208",
    "locale": "zh",
    "title": "樱花树下害羞双马尾少女",
    "category": "场景与叙事",
    "description": "[中文] 生成一张高质量二次元美少女图片。 角色设定： - 年龄：17岁 - 发型：双马尾，颜色：樱花粉，发梢带点渐变紫色 - 眼睛：大而明亮，紫色瞳孔，有星星高光 - 服装：JK制服，白色衬衫，深蓝色格子裙，红色领结 - 配饰：白色过膝袜，棕色小皮鞋，头上戴一个粉色蝴蝶结 风格",
    "promptPreview": "[中文] 生成一张高质量二次元美少女图片。 角色设定： - 年龄：17岁 - 发型：双马尾，颜色：樱花粉，发梢带点渐变紫色 - 眼睛：大而明亮，紫色瞳孔，有星星高光 - 服装：JK制服，白色衬衫，深蓝色格子裙，红色领结 - 配饰：白色过膝袜，棕色小皮鞋，头上戴一个粉色蝴蝶结 风格要求： - 日系动画风格，线条清晰 - ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@joshesye",
    "tags": [
      "Scenes & Storytelling",
      "UI",
      "Character",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-329",
    "locale": "zh",
    "title": "烬甲猎鹰者与燃翼神禽",
    "category": "场景与叙事",
    "description": "[中文] 一幅充满奇幻色彩的电影场景：一位英姿飒爽的女战士兼猎鹰师，身着饱经战火洗礼、饰以闪耀余烬纹理的皮甲，漫步于幽暗迷雾笼罩的森林之中。她高举手臂，指挥着一头巨大的凤凰与雄鹰的混合体，这头猛禽双翼燃烧，羽毛燃焰，尖端喷吐着火焰。它周身散发着橙红色的熔岩光芒，火星和余烬飞溅。女",
    "promptPreview": "[中文] 一幅充满奇幻色彩的电影场景：一位英姿飒爽的女战士兼猎鹰师，身着饱经战火洗礼、饰以闪耀余烬纹理的皮甲，漫步于幽暗迷雾笼罩的森林之中。她高举手臂，指挥着一头巨大的凤凰与雄鹰的混合体，这头猛禽双翼燃烧，羽毛燃焰，尖端喷吐着火焰。它周身散发着橙红色的熔岩光芒，火星和余烬飞溅。女战士梳着辫子，皮肤上沾满了灰烬，神情坚定",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@iamsofiaijaz",
    "tags": [
      "Scenes & Storytelling",
      "Realistic",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-390",
    "locale": "en",
    "title": "羊毛毡国家微缩世界",
    "category": "场景与叙事",
    "description": "Country: [INSERT COUNTRY NAME] A miniature felt diorama made of fluffy yarn, wool, and needlework, designed as a single cohesive small world",
    "promptPreview": "Country: [INSERT COUNTRY NAME] A miniature felt diorama made of fluffy yarn, wool, and needlework, designed as a single cohesive small world (not a collage), re",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@volkan_iras",
    "tags": [
      "Scenes & Storytelling",
      "UI",
      "Realistic",
      "Brand",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-395",
    "locale": "en",
    "title": "骑士法师大战石像魔像",
    "category": "场景与叙事",
    "description": "Create a cinematic dark fantasy action scene in a ruined cathedral hall: a {argument name=\"hero type\" default=\"female armored knight-mage\"} ",
    "promptPreview": "Create a cinematic dark fantasy action scene in a ruined cathedral hall: a {argument name=\"hero type\" default=\"female armored knight-mage\"} crouches in a defens",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@RamonVi25791296",
    "tags": [
      "Scenes & Storytelling",
      "UI",
      "Realistic",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-319",
    "locale": "zh",
    "title": "鸟群织就的梦幻高定时装秀",
    "category": "场景与叙事",
    "description": "[中文] 一个充满趣味的高级时装T台场景，主角是一位自信的女性，正走在奢华时装秀的T台上，身穿一件完全由鸟类制成的非凡高级定制礼服。数百只优雅、色彩鲜艳的鸟类构成了飘逸的雕塑感礼服形状，像活着的羽毛一样层叠，翅膀微微张开，营造出布料和运动的错觉。一些鸟儿在她周围轻轻升入空中，捕捉",
    "promptPreview": "[中文] 一个充满趣味的高级时装T台场景，主角是一位自信的女性，正走在奢华时装秀的T台上，身穿一件完全由鸟类制成的非凡高级定制礼服。数百只优雅、色彩鲜艳的鸟类构成了飘逸的雕塑感礼服形状，像活着的羽毛一样层叠，翅膀微微张开，营造出布料和运动的错觉。一些鸟儿在她周围轻轻升入空中，捕捉于飞行瞬间，增添了神奇、超现实的运动感。",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrDasOnX",
    "tags": [
      "Scenes & Storytelling",
      "Character",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-2",
    "locale": "en",
    "title": "UI与界面 - JSON 进阶模板（推荐给 Agent 调用）",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 UI与界面 模板。",
    "promptPreview": "{ \"type\": \"UI Screenshot\", \"platform\": \"iOS\", \"product\": \"Fitness App\", \"layout\": \"Card-based feed with bottom tab bar\", \"style\": { \"theme\": \"Dark Mode\", \"prima",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "UI与界面",
      "json"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-1",
    "locale": "zh",
    "title": "UI与界面 - 常规模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 UI与界面 模板。",
    "promptPreview": "为[产品类型]生成一张[平台，如 iOS/Android/Web]界面图。 核心功能：[功能点A]、[功能点B]、[功能点C]。 视觉风格：[极简/科技/拟物]，主色[颜色]，强调色[颜色]。 布局：[顶部导航/双栏/卡片流]，信息层级清晰，留白充足。 输出：高保真UI截图，文字清晰可读，比例[9:16/16:9]。",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "UI与界面",
      "text"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-3",
    "locale": "zh",
    "title": "UI与界面 - 截图生成模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 UI与界面 模板。",
    "promptPreview": "生成一张[平台，如 X/抖音/小红书/微信朋友圈]内容截图，[深色/浅色]模式。 整体比例：[9:16 / 3:4 / 1:1]，手机截图风格。 核心内容： - 账号信息：[头像描述 / 用户名 / 认证标识] - 正文内容：[具体文本内容，包含指定中文] - 互动数据：[点赞/评论/转发/收藏数量] 界面元素： -",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "UI与界面",
      "text"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-4",
    "locale": "zh",
    "title": "UI与界面 - 直播界面模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 UI与界面 模板。",
    "promptPreview": "生成一张[平台，如抖音/快手/B站]直播界面截图。 主播：[人物描述/名称]，姿态：[坐姿/站立/动作]，服装：[服装描述]。 背景：[直播间背景描述]，灯光：[暖色/冷色/混合]。 UI叠加层： - 顶部：主播头像 + 关注按钮 + 在线人数 + 排名/热值 - 左下：弹幕/评论列表（[N]条，内容示例） - 右下或",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "UI与界面",
      "text"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-6",
    "locale": "en",
    "title": "图表与信息可视化 - JSON 进阶模板（推荐给 Agent 调用）",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 图表与信息可视化 模板。",
    "promptPreview": "{ \"type\": \"Infographic\", \"topic\": \"Urban Metabolism\", \"audience\": \"General Public\", \"structure\": { \"title_area\": \"城市生命系统图谱\", \"layout\": \"Isometric cutaway, 12 nu",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "图表与信息可视化",
      "json"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-7",
    "locale": "zh",
    "title": "图表与信息可视化 - 尺度缩放科学信息图模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 图表与信息可视化 模板。",
    "promptPreview": "为[主题]生成一张科学尺度缩放信息图。 结构：6-8 个圆形或六边形框，按从微观到宏观的尺度递进排列。 每个框包含：尺度名称、3-5 个词的洞察、测量单位或放大倍率，以及该尺度下的高细节 3D 渲染。 用细线连接各尺度，避免重复层级。标题使用“[主题]：AT EVERY SCALE”或“ZOOM: THE WORLD",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "图表与信息可视化",
      "text"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-5",
    "locale": "zh",
    "title": "图表与信息可视化 - 常规模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 图表与信息可视化 模板。",
    "promptPreview": "生成[主题：明确、具体，避免宽泛。例如：“老年人日常健康管理指南”而非“健康”]信息图，目标读者为[人群:细化人群特征，如年龄段、职业、兴趣等]。 结构：标题区 + [3-5]个模块（每模块含图标、短标题、1-2句说明,模块间逻辑：可用箭头、颜色区分或连接线提示信息流或关系等适当的方式）。 图表类型：[流程图/对比图/",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "图表与信息可视化",
      "text"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-12",
    "locale": "zh",
    "title": "海报与排版 - 单款签名提取模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 海报与排版 模板。",
    "promptPreview": "从输入图中的[位置/编号/风格名]签名里，提取该签名的核心笔势，生成一张纯签名图。 要求： - 只保留签名主体，不生成海报卡片、标题、副标题或说明文字 - 保留原签名的起笔、连笔、结构倾斜、飞白和收笔节奏 - 背景为纯白或极浅米白 - 签名居中，尺寸充足，边缘留白干净 - 墨色为深黑或墨黑，带自然笔锋、轻微墨痕和真实手",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "海报与排版",
      "text"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-11",
    "locale": "zh",
    "title": "海报与排版 - 多风格签名选择海报模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 海报与排版 模板。",
    "promptPreview": "你是一个高端签名设计系统 + 风格人格视觉系统。 输入： 姓名：[姓名/昵称] 任务： 基于姓名自动生成一张 9:16 竖版「多风格签名选择海报」。 目标是把姓名转译成 6 种具有笔势、气质和力量感的签名方案。 隐藏分析： 1. 分析姓名字形：疏密、横竖比例、重心、连笔空间、草写空间。 2. 推断气质：清冷、张扬、克制",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "海报与排版",
      "text"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-8",
    "locale": "zh",
    "title": "海报与排版 - 常规模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 海报与排版 模板。",
    "promptPreview": "设计一张[活动/产品/电影]海报，主题为[主题词]。 主视觉：[主体元素]，标题文案：[标题]，副标题：[副标题]。 版式：[居中/左对齐/对角构图]，风格：[复古/未来/极简]。 色彩：[主色 + 辅色]，氛围：[情绪关键词]。 输出：可用于社媒传播的高分辨率海报。",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "海报与排版",
      "text"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-10",
    "locale": "en",
    "title": "海报与排版 - 概念字体海报模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 海报与排版 模板。",
    "promptPreview": "Create ONE finished premium conceptual typography poster for the exact title: \"[标题/词语/短句]\" Single poster only. No moodboard, grid, presentation board, mockup, c",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "海报与排版",
      "text"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-13",
    "locale": "zh",
    "title": "海报与排版 - 签名练习拆解图模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 海报与排版 模板。",
    "promptPreview": "基于输入的[签名图/签名风格]，生成一张签名练习拆解图。 目标： 帮助用户用黑笔在纸上练好这个签名，拆解每一笔的书写路径、顺序、力度和节奏。 画面结构： - 竖版教学图或横版练习板 - 顶部放最终签名小样 - 中部用 8-12 个步骤拆解关键笔画 - 每个步骤展示当前笔画、运动方向箭头、起笔点、停顿点、收笔点 - 下方",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "海报与排版",
      "text"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-14",
    "locale": "zh",
    "title": "海报与排版 - 签名练习拆解图模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 海报与排版 模板。",
    "promptPreview": "**避坑指南** - **不要偷懒**：写清“主视觉到底是什么玩意儿”，不要只丢一句“做一张海报”就指望出神图。 - **文案硬编码**：主标题与副标题都要写死，否则模型会给你疯狂加戏，自动瞎编不知所云的文字。 - **运动海报先定结构**：运动 Campaign 最容易变成杂乱拼贴，先锁定“单主视觉 / 三联画 /",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "海报与排版"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-15",
    "locale": "zh",
    "title": "海报与排版 - 签名练习拆解图模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 海报与排版 模板。",
    "promptPreview": "**避坑指南** - **材质和光影是灵魂**：一定要堆叠材质（如“磨砂质感”）和灯光（如“轮廓光”）的关键词，商品图一旦没有光影，立刻变成地摊货。 - **别把促销贴满全屏**：文案只给核心的 1-2 句（如“新品上市”），字多了画面就毁了。 - **先分析再出图**：美妆推荐类不要直接让模型摆色号，先要求它分析肤色",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "海报与排版"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-16",
    "locale": "zh",
    "title": "海报与排版 - 签名练习拆解图模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 海报与排版 模板。",
    "promptPreview": "**避坑指南** - **做减法**：先定义品牌关键词，再要求视觉输出，结果更统一。别让它画“一条喷火的龙缠绕在长城的柱子上还带着闪电”，那不叫 Logo，那叫插画。 - **强制背景**：必须强调“纯白背景（Pure White Background）”，方便后期抠图。 - **先做品牌战略再画 Logo**：如果缺",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "海报与排版"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-17",
    "locale": "zh",
    "title": "海报与排版 - 签名练习拆解图模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 海报与排版 模板。",
    "promptPreview": "**避坑指南** - **控制视角**：建筑图最容易翻车的就是透视变形。用“Eye-level perspective（人眼视角）”能压住它。 - **冷暖对比**：室外的冷光（蓝/灰）和室内的暖光（黄/橙）搭配，是提升空间高级感的作弊码。 <a name=\"tpl-photo\"></a> ### 摄影与写实 **常规",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "海报与排版"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-18",
    "locale": "zh",
    "title": "海报与排版 - 签名练习拆解图模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 海报与排版 模板。",
    "promptPreview": "**避坑指南** - **加点瑕疵**：AI 画的人太完美了，反而像假人。加入“皮肤纹理（skin pores）”、“雀斑”、“轻微胶片颗粒（film grain）”，真实感瞬间拉满。 - **用参数说话**：用 `f/1.4` 代替“浅景深”，用 `50mm` 代替“半身照”，大模型吃这套。 - **把“不完美”写具",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "海报与排版"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-19",
    "locale": "zh",
    "title": "海报与排版 - 签名练习拆解图模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 海报与排版 模板。",
    "promptPreview": "**避坑指南** - **锁定笔触**：插画如果不限制笔触（如“厚涂”、“水彩晕染”），它通常会给你一种毫无灵魂的 AI 默认塑料风。 - **慎用大师名**：提大师名字很爽，但容易被模型原样照搬其代表作的构图。建议提取大师的特征（如“梵高的旋转星空笔触”），而不是直接写大师名。 <a name=\"tpl-charac",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "海报与排版"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-20",
    "locale": "zh",
    "title": "海报与排版 - 签名练习拆解图模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 海报与排版 模板。",
    "promptPreview": "**避坑指南** - **拆解五官**：不要只写“很美的女孩”，大模型不知道你的审美标准。拆解成“桃花眼、高鼻梁、野生眉”。 - **服装材质**：写清衣服的材质（如“丝绸”、“机能防风面料”），能让角色立刻变得立体。 - **动作表要锁网格**：动作分解图必须明确面板数量、编号、每格结构，否则模型会把步骤挤成一张杂乱",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "海报与排版"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-21",
    "locale": "zh",
    "title": "海报与排版 - 签名练习拆解图模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 海报与排版 模板。",
    "promptPreview": "**避坑指南** - **要有“动词”**：叙事图最怕画成风景明信片。一定要写“事件”（如“正在崩塌”、“刚点燃火把”），让画面动起来。 - **镜头语言**：使用“Low angle shot（低角度仰拍）”或“Dutch angle（倾斜镜头）”来增加戏剧冲突。 <a name=\"tpl-history\"></a>",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "海报与排版"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-22",
    "locale": "zh",
    "title": "海报与排版 - 签名练习拆解图模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 海报与排版 模板。",
    "promptPreview": "**避坑指南** - **拒绝大杂烩**：明确朝代（唐/宋/明），否则大模型会给你画出一个穿着和服、拿着清朝折扇在唐朝宫殿里的人。 - **强制排雷**：一定要加上“禁用现代元素（No modern elements）”，防止古风美女手里突然多出一杯星巴克。 <a name=\"tpl-document\"></a> ##",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "海报与排版"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-23",
    "locale": "zh",
    "title": "海报与排版 - 签名练习拆解图模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 海报与排版 模板。",
    "promptPreview": "**企业画册系统模板** > 来源参考：[@MrLarus](https://x.com/MrLarus/status/2056974720893939950)",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "海报与排版"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-24",
    "locale": "zh",
    "title": "海报与排版 - 签名练习拆解图模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 海报与排版 模板。",
    "promptPreview": "**避坑指南** - **结构优先**：明确“栏数（columns）”和“留白（margins）”比堆砌风格词更重要。 - **放弃全文**：不要指望大模型能排出一整页毫无错字的正文，让它用“模拟文本（Simulated text blocks）”填充正文，只写死大标题。 - **系统预览**：企业画册类任务最好补一张",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "海报与排版"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-25",
    "locale": "zh",
    "title": "海报与排版 - 签名练习拆解图模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 海报与排版 模板。",
    "promptPreview": "**避坑指南** - **先说干嘛的**：一上来先写“任务目标和用途”，让模型建立全局上下文，再写视觉细节。 - **A/B 测试**：通用场景建议在 prompt 里要求“一次生成主方案 + 备选方案”，方便你直接挑好的。 ***",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "海报与排版"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-template-9",
    "locale": "zh",
    "title": "海报与排版 - 运动商业 Campaign 模板",
    "category": "工业级模板",
    "description": "来自 freestylefly 工业级提示词模板与防坑指南的 海报与排版 模板。",
    "promptPreview": "设计一张[运动项目/健身品类]商业 Campaign 海报。 主体：[运动员/模特/产品道具]，姿态：[坐姿/冲刺/挥拍/力量动作]。 核心道具：[球拍/哑铃/球鞋/球衣]，以夸张比例或对角构图成为视觉锚点。 版式：[单张强主视觉/三联画/数据涂鸦海报]。 大字标题：\"[主标题]\"，辅助文案：\"[短句/数据/精神口号]",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "freestylefly templates",
    "tags": [
      "template",
      "海报与排版",
      "text"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-381",
    "locale": "en",
    "title": "90 年代公寓场景参考板",
    "category": "建筑与空间",
    "description": "{ \"type\": \"scene reference board — 90s apartment living room, cinematic night\", \"style\": \"cinematic film photography, 35mm grain, warm amber",
    "promptPreview": "{ \"type\": \"scene reference board — 90s apartment living room, cinematic night\", \"style\": \"cinematic film photography, 35mm grain, warm amber shadow fill, deep c",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Iancu_ai",
    "tags": [
      "Architecture & Spaces",
      "Poster",
      "Realistic",
      "Product",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-11",
    "locale": "zh",
    "title": "一张手绘风格的城市美食地图，以台州为主题",
    "category": "建筑与空间",
    "description": "一张手绘风格的城市美食地图，以台州为主题。画面以鸟瞰视角的手绘简化城市地图为底，标注椒江、路桥、黄岩等区域和灵江、台州湾等水系地标，不追求精确比例而是追求可爱的水彩手绘感。地图上分布着12个美食地点的精致手绘小插画：1. 椒江老粮坊的蛋清羊尾（金黄蓬松的蛋泡甜点撒着糖粉，筷子夹起",
    "promptPreview": "一张手绘风格的城市美食地图，以台州为主题。画面以鸟瞰视角的手绘简化城市地图为底，标注椒江、路桥、黄岩等区域和灵江、台州湾等水系地标，不追求精确比例而是追求可爱的水彩手绘感。地图上分布着12个美食地点的精致手绘小插画：1. 椒江老粮坊的蛋清羊尾（金黄蓬松的蛋泡甜点撒着糖粉，筷子夹起拉丝）2. 临海紫阳古街的食饼筒（一个饱",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "小红书号510244722",
    "tags": [
      "Architecture & Spaces",
      "Illustration",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-489",
    "locale": "en",
    "title": "城市地图微缩旅行海报",
    "category": "建筑与空间",
    "description": "Create a highly detailed cinematic miniature tilt-shift travel scene of [CITY NAME] featuring a realistic [VEHICLE NAME] driving along a win",
    "promptPreview": "Create a highly detailed cinematic miniature tilt-shift travel scene of [CITY NAME] featuring a realistic [VEHICLE NAME] driving along a winding elevated road t",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@SimplyAnnisa",
    "tags": [
      "Architecture & Spaces",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-211",
    "locale": "zh",
    "title": "天坛古建拆解全图",
    "category": "建筑与空间",
    "description": "[中文] 生成一个天坛的建筑拆解图，有详细的说明，中式美学风格 [English] Generate an architectural exploded view of the Temple of Heaven, with detailed annotations, Chines",
    "promptPreview": "[中文] 生成一个天坛的建筑拆解图，有详细的说明，中式美学风格 [English] Generate an architectural exploded view of the Temple of Heaven, with detailed annotations, Chinese aesthetic style",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@TanShilong",
    "tags": [
      "Architecture & Spaces",
      "Architecture",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-50",
    "locale": "en",
    "title": "建筑空间场景图",
    "category": "建筑与空间",
    "description": "A highly detailed, cinematic wide shot of a grand, dark gothic hall with a {argument name=\"atmosphere\" default=\"dark fantasy\"} aesthetic. In",
    "promptPreview": "A highly detailed, cinematic wide shot of a grand, dark gothic hall with a {argument name=\"atmosphere\" default=\"dark fantasy\"} aesthetic. In the center, a singl",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@nomen_machine",
    "tags": [
      "Architecture & Spaces",
      "Architecture",
      "Tech",
      "Fashion",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-274",
    "locale": "zh",
    "title": "成都吃货暴走手绘美食地图",
    "category": "建筑与空间",
    "description": "[中文] 一张手绘风格的城市美食地图，以成都为主题。画面以鸟瞰视角的手绘简化城市地图为底，标注主要道路和地标但不追求精确比例而是追求可爱的手绘感。地图上分布着 12 个美食地点的精致手绘小插画：春熙路的串串香（一把竹签插着各种食材冒着热气）、宽窄巷子的三大炮（三个糯米团子飞向铜盘",
    "promptPreview": "[中文] 一张手绘风格的城市美食地图，以成都为主题。画面以鸟瞰视角的手绘简化城市地图为底，标注主要道路和地标但不追求精确比例而是追求可爱的手绘感。地图上分布着 12 个美食地点的精致手绘小插画：春熙路的串串香（一把竹签插着各种食材冒着热气）、宽窄巷子的三大炮（三个糯米团子飞向铜盘）、建设路的蛋烘糕（金黄酥脆正在翻面）、",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@Panda20230902",
    "tags": [
      "Architecture & Spaces",
      "UI",
      "Illustration",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-369",
    "locale": "zh",
    "title": "明洞旅游区域地图",
    "category": "建筑与空间",
    "description": "[エリア]の観光エリアマップを画像で作成して",
    "promptPreview": "[エリア]の観光エリアマップを画像で作成して",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@so_ainsight",
    "tags": [
      "Architecture & Spaces",
      "Architecture",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-411",
    "locale": "en",
    "title": "极简建筑地标海报",
    "category": "建筑与空间",
    "description": "Design a luxury minimalist poster centered on a famous architectural landmark of your choice ([building name]). The focal element is an illu",
    "promptPreview": "Design a luxury minimalist poster centered on a famous architectural landmark of your choice ([building name]). The focal element is an illustrated rendering of",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@iamaiistudio",
    "tags": [
      "Architecture & Spaces",
      "UI",
      "Poster",
      "Brand",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-331",
    "locale": "zh",
    "title": "西安手绘水彩城市地图",
    "category": "建筑与空间",
    "description": "生成一张手绘水彩风格的「西安」城市地图，包含当地特色美食、地标建筑及城市特色",
    "promptPreview": "生成一张手绘水彩风格的「西安」城市地图，包含当地特色美食、地标建筑及城市特色",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "苍何原创实测（公众号文章《我逆向了 329 条 GPT-Image2 提示词模板，全部开源！》）",
    "tags": [
      "Architecture & Spaces",
      "Architecture",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-117",
    "locale": "en",
    "title": "主题海报版式设计",
    "category": "插画与艺术",
    "description": "A humorous 3D cartoon illustration of a therapy session in a cozy office. On the left, a {argument name=\"patient character\" default=\"sad ant",
    "promptPreview": "A humorous 3D cartoon illustration of a therapy session in a cozy office. On the left, a {argument name=\"patient character\" default=\"sad anthropomorphic avocado",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@nicdunz",
    "tags": [
      "Illustration & Art",
      "Poster",
      "Illustration",
      "Character",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-446",
    "locale": "en",
    "title": "低多边形纸艺男士肖像",
    "category": "插画与艺术",
    "description": "Create a highly detailed low-poly papercraft portrait of a stylish young man, designed like an origami paper sculpture. The character has cu",
    "promptPreview": "Create a highly detailed low-poly papercraft portrait of a stylish young man, designed like an origami paper sculpture. The character has curly dark brown hair,",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@iamsofiaijaz",
    "tags": [
      "Illustration & Art",
      "Realistic",
      "Character",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-112",
    "locale": "en",
    "title": "信息图可视化设计",
    "category": "插画与艺术",
    "description": "Generate a 12-grid card image of the 12 Golden Saints from Saint Seiya, with each card featuring its corresponding Chinese name, 4 cards per",
    "promptPreview": "Generate a 12-grid card image of the 12 Golden Saints from Saint Seiya, with each card featuring its corresponding Chinese name, 4 cards per row, in a 16:9 aspe",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@songguoxiansen",
    "tags": [
      "Illustration & Art",
      "Infographic",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-52",
    "locale": "en",
    "title": "写实摄影风格图",
    "category": "插画与艺术",
    "description": "A realistic photograph of a whiteboard with a highly detailed {argument name=\"marker color\" default=\"green\"} dry-erase marker drawing of {ar",
    "promptPreview": "A realistic photograph of a whiteboard with a highly detailed {argument name=\"marker color\" default=\"green\"} dry-erase marker drawing of {argument name=\"subject",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@nicdunz",
    "tags": [
      "Illustration & Art",
      "Realistic",
      "Character",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-30",
    "locale": "en",
    "title": "写实摄影风格图",
    "category": "插画与艺术",
    "description": "Express [{argument name=\"subject\" default=\"a powerful AI builder\"}] in a graffiti sketch style, presenting an overall visual effect of rapid",
    "promptPreview": "Express [{argument name=\"subject\" default=\"a powerful AI builder\"}] in a graffiti sketch style, presenting an overall visual effect of rapid sketching, free tra",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@opc_8838",
    "tags": [
      "Illustration & Art",
      "UI",
      "Realistic",
      "Illustration",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-316",
    "locale": "zh",
    "title": "冲破次元壁的写实漫画跑者",
    "category": "插画与艺术",
    "description": "[中文] { \"prompt\": \"超写实，一位留着深色短卷发、修剪整齐的胡须和黑色方形眼镜的年轻男子的鲜艳逼真渲染，身穿深色纹理高领毛衣和牛仔裤。他奔跑到一半被捕捉下来，姿态充满动感，向前突破，充满戏剧性地从一个破碎的漫画分镜框中显现——一条腿和一只手臂冲入现实世界，而身体的其",
    "promptPreview": "[中文] { \"prompt\": \"超写实，一位留着深色短卷发、修剪整齐的胡须和黑色方形眼镜的年轻男子的鲜艳逼真渲染，身穿深色纹理高领毛衣和牛仔裤。他奔跑到一半被捕捉下来，姿态充满动感，向前突破，充满戏剧性地从一个破碎的漫画分镜框中显现——一条腿和一只手臂冲入现实世界，而身体的其余部分仍留在漫画框内。他的表情充满活力和",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@Fujimoto_hina",
    "tags": [
      "Illustration & Art",
      "Realistic",
      "Character",
      "3D",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-113",
    "locale": "en",
    "title": "动漫插画创作图",
    "category": "插画与艺术",
    "description": "A highly detailed anime illustration of a fierce female warrior with long flowing {argument name=\"hair color\" default=\"black\"} hair and pier",
    "promptPreview": "A highly detailed anime illustration of a fierce female warrior with long flowing {argument name=\"hair color\" default=\"black\"} hair and piercing {argument name=",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@REd8358",
    "tags": [
      "Illustration & Art",
      "UI",
      "Illustration",
      "Tech",
      "Social",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-193",
    "locale": "zh",
    "title": "千手观音化身打工人",
    "category": "插画与艺术",
    "description": "[中文] 一幅高度详细的千手观音菩萨工笔画。 然而，千手并没有拿着神圣的宗教法器，而是拿着现代办公和家用物品：**笔记本电脑、智能手机、成堆的文件、咖啡杯、印章、计算器、拖把和奶瓶**。它代表了终极的多任务处理现代工作者。 脑后的金色光环由旋转的时钟齿轮组成。 **在右下角，一个",
    "promptPreview": "[中文] 一幅高度详细的千手观音菩萨工笔画。 然而，千手并没有拿着神圣的宗教法器，而是拿着现代办公和家用物品：**笔记本电脑、智能手机、成堆的文件、咖啡杯、印章、计算器、拖把和奶瓶**。它代表了终极的多任务处理现代工作者。 脑后的金色光环由旋转的时钟齿轮组成。 **在右下角，一个单一的红色竖排艺术家印章写着“吴先生”（",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@johnAGI168",
    "tags": [
      "Illustration & Art",
      "Illustration",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-405",
    "locale": "en",
    "title": "可爱纸艺风照片重绘",
    "category": "插画与艺术",
    "description": "Recreate this image in a paper craft style, simplifying the details to make them suitable for paper craft artwork. Arrange the overall compo",
    "promptPreview": "Recreate this image in a paper craft style, simplifying the details to make them suitable for paper craft artwork. Arrange the overall composition to feel visua",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@oggii_0",
    "tags": [
      "Illustration & Art",
      "UI",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-115",
    "locale": "en",
    "title": "品牌视觉识别图",
    "category": "插画与艺术",
    "description": "{ \"type\": \"two-page manga spread\", \"style\": \"highly detailed realistic manga, monochrome, screentones, dramatic lighting, psychological thri",
    "promptPreview": "{ \"type\": \"two-page manga spread\", \"style\": \"highly detailed realistic manga, monochrome, screentones, dramatic lighting, psychological thriller\", \"global_eleme",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@onofumi_AI",
    "tags": [
      "Illustration & Art",
      "UI",
      "Realistic",
      "Brand",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-410",
    "locale": "en",
    "title": "夸张动漫风主体重绘",
    "category": "插画与艺术",
    "description": "Create a trending anime art style image from the uploaded subject. Use confident line-work with slight variation and minimal cel shading usi",
    "promptPreview": "Create a trending anime art style image from the uploaded subject. Use confident line-work with slight variation and minimal cel shading using flat shadow shape",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Zyrellix",
    "tags": [
      "Illustration & Art",
      "Realistic",
      "Character",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-461",
    "locale": "en",
    "title": "家庭旅行纸雕拼贴",
    "category": "插画与艺术",
    "description": "Adorable kawaii family travel collage, combining four different scenes in one composition, cute anime-inspired parents, friends and children",
    "promptPreview": "Adorable kawaii family travel collage, combining four different scenes in one composition, cute anime-inspired parents, friends and children taking selfies and ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Taaruk_",
    "tags": [
      "Illustration & Art",
      "Illustration",
      "Character",
      "3D",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-435",
    "locale": "en",
    "title": "层叠纸雕情侣插画",
    "category": "插画与艺术",
    "description": "{ \"style\": \"layered paper-cut illustration, papercraft diorama, handcrafted aesthetic\", \"technique\": { \"layering\": \"multiple stacked paper l",
    "promptPreview": "{ \"style\": \"layered paper-cut illustration, papercraft diorama, handcrafted aesthetic\", \"technique\": { \"layering\": \"multiple stacked paper layers with soft drop",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Just_sharon7",
    "tags": [
      "Illustration & Art",
      "Realistic",
      "Illustration",
      "Character",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-458",
    "locale": "en",
    "title": "巴黎秋季时装插画",
    "category": "插画与艺术",
    "description": "Full-body fashion illustration of a young woman walking down a Parisian street in autumn, wearing a long camel wool coat draped open over a ",
    "promptPreview": "Full-body fashion illustration of a young woman walking down a Parisian street in autumn, wearing a long camel wool coat draped open over a black ribbed turtlen",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@heyfatema",
    "tags": [
      "Illustration & Art",
      "UI",
      "Realistic",
      "Illustration",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-495",
    "locale": "en",
    "title": "巴黎街头故事书插画",
    "category": "插画与艺术",
    "description": "Portrait illustration in a storybook style featuring a young adult woman exploring the streets of Paris. She is laughing happily with her ey",
    "promptPreview": "Portrait illustration in a storybook style featuring a young adult woman exploring the streets of Paris. She is laughing happily with her eyes closed, holding a",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@MissDelulu9",
    "tags": [
      "Illustration & Art",
      "Illustration",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-127",
    "locale": "en",
    "title": "建筑空间场景图",
    "category": "插画与艺术",
    "description": "{ \"type\": \"manga page\", \"style\": \"anime illustration, full color\", \"characters\": { \"woman\": { \"appearance\": \"long {argument name=\\\"hair colo",
    "promptPreview": "{ \"type\": \"manga page\", \"style\": \"anime illustration, full color\", \"characters\": { \"woman\": { \"appearance\": \"long {argument name=\\\"hair color\\\" default=\\\"black\\",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@studiomasakaki",
    "tags": [
      "Illustration & Art",
      "Illustration",
      "Character",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-121",
    "locale": "en",
    "title": "建筑空间场景图",
    "category": "插画与艺术",
    "description": "{ \"type\": \"3-panel manga page\", \"style\": \"anime, highly detailed, cinematic lighting, futuristic corporate\", \"layout\": { \"structure\": \"1 wid",
    "promptPreview": "{ \"type\": \"3-panel manga page\", \"style\": \"anime, highly detailed, cinematic lighting, futuristic corporate\", \"layout\": { \"structure\": \"1 wide top panel, 2 squar",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@loilokoji",
    "tags": [
      "Illustration & Art",
      "UI",
      "Character",
      "Tech",
      "Travel",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-120",
    "locale": "en",
    "title": "建筑空间场景图",
    "category": "插画与艺术",
    "description": "A dynamic anime illustration of a girl with spiky {argument name=\"hair color\" default=\"blonde\"} hair tied in a high ponytail with a black bo",
    "promptPreview": "A dynamic anime illustration of a girl with spiky {argument name=\"hair color\" default=\"blonde\"} hair tied in a high ponytail with a black bow, striking teal eye",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@UNIBRACITY",
    "tags": [
      "Illustration & Art",
      "Illustration",
      "Character",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-374",
    "locale": "zh",
    "title": "彩色潦草小狗线条风格重绘",
    "category": "插画与艺术",
    "description": "彩色潦草小狗线条风格绘制该图，童趣和doodle加入其中，务必使用毫无章法的绘制手法，凌乱和草率即可。",
    "promptPreview": "彩色潦草小狗线条风格绘制该图，童趣和doodle加入其中，务必使用毫无章法的绘制手法，凌乱和草率即可。",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@berryxia",
    "tags": [
      "Illustration & Art",
      "Illustration",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-409",
    "locale": "en",
    "title": "拙劣 MS Paint 风重绘",
    "category": "插画与艺术",
    "description": "Please redraw the attached image in the most clumsy, messy, and hopelessly pathetic way possible. Use a white background and make it look li",
    "promptPreview": "Please redraw the attached image in the most clumsy, messy, and hopelessly pathetic way possible. Use a white background and make it look like it was drawn in M",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Ciri_ai",
    "tags": [
      "Illustration & Art",
      "Illustration",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-114",
    "locale": "en",
    "title": "插画艺术创作图",
    "category": "插画与艺术",
    "description": "{ \"type\": \"4-panel vertical comic strip\", \"style\": \"{argument name=\\\"art style\\\" default=\\\"black and white pencil sketch, crosshatch shading",
    "promptPreview": "{ \"type\": \"4-panel vertical comic strip\", \"style\": \"{argument name=\\\"art style\\\" default=\\\"black and white pencil sketch, crosshatch shading, satirical caricatu",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@kaikaitheaiguy",
    "tags": [
      "Illustration & Art",
      "UI",
      "Illustration",
      "Character",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-43",
    "locale": "en",
    "title": "插画艺术创作图",
    "category": "插画与艺术",
    "description": "{ \"type\": \"character portrait grid\", \"theme\": \"{argument name=\\\"theme\\\" default=\\\"Game of Thrones characters\\\"}\", \"style\": \"{argument name=\\",
    "promptPreview": "{ \"type\": \"character portrait grid\", \"theme\": \"{argument name=\\\"theme\\\" default=\\\"Game of Thrones characters\\\"}\", \"style\": \"{argument name=\\\"art style\\\" default",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@stark_nico99",
    "tags": [
      "Illustration & Art",
      "Illustration",
      "Character",
      "Tech",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-34",
    "locale": "en",
    "title": "插画艺术创作图",
    "category": "插画与艺术",
    "description": "{ \"type\": \"character avatar grid\", \"theme\": \"{argument name=\\\"theme\\\" default=\\\"Journey to the West mythology\\\"}\", \"style\": \"{argument name=",
    "promptPreview": "{ \"type\": \"character avatar grid\", \"theme\": \"{argument name=\\\"theme\\\" default=\\\"Journey to the West mythology\\\"}\", \"style\": \"{argument name=\\\"art style\\\" defaul",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Ryan_Suo",
    "tags": [
      "Illustration & Art",
      "Illustration",
      "Character",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-32",
    "locale": "en",
    "title": "插画艺术创作图",
    "category": "插画与艺术",
    "description": "{ \"type\": \"3x3 character expression grid\", \"style\": \"{argument name=\\\"art style\\\" default=\\\"3D animation, Pixar style\\\"}\", \"character_base\":",
    "promptPreview": "{ \"type\": \"3x3 character expression grid\", \"style\": \"{argument name=\\\"art style\\\" default=\\\"3D animation, Pixar style\\\"}\", \"character_base\": \"{argument name=\\\"c",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@austinit",
    "tags": [
      "Illustration & Art",
      "Illustration",
      "Character",
      "3D",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-445",
    "locale": "en",
    "title": "旅游照水墨明信片",
    "category": "插画与艺术",
    "description": "Create a dreamy watercolor travel illustration style from the attached photo. Use hand-painted urban sketchbook aesthetic, delicate ink line",
    "promptPreview": "Create a dreamy watercolor travel illustration style from the attached photo. Use hand-painted urban sketchbook aesthetic, delicate ink linework mixed with soft",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@zhgqthomas",
    "tags": [
      "Illustration & Art",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-423",
    "locale": "en",
    "title": "日系手绘涂鸦半身插画",
    "category": "插画与艺术",
    "description": "Generate an illustration of \"me\" as you imagine it. Features include a Japanese illustration style, distinct character features, natural emo",
    "promptPreview": "Generate an illustration of \"me\" as you imagine it. Features include a Japanese illustration style, distinct character features, natural emotional expressions, ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@heyfatema",
    "tags": [
      "Illustration & Art",
      "UI",
      "Illustration",
      "Character",
      "Tech",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-476",
    "locale": "en",
    "title": "早安拿铁微缩女孩",
    "category": "插画与艺术",
    "description": "Create ultra-fine highly detailed 3D realistic miniature chibi-like cute girl, wearing cream colour top and jeans, resting and floating on c",
    "promptPreview": "Create ultra-fine highly detailed 3D realistic miniature chibi-like cute girl, wearing cream colour top and jeans, resting and floating on creamy latte cup, sty",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Zyrellix",
    "tags": [
      "Illustration & Art",
      "Realistic",
      "3D",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-479",
    "locale": "en",
    "title": "杂志纸艺拼贴重绘",
    "category": "插画与艺术",
    "description": "Transform the uploaded image into a minimalist illustration in a magazine collage style, using paper cutouts. Retain the main subject, pose,",
    "promptPreview": "Transform the uploaded image into a minimalist illustration in a magazine collage style, using paper cutouts. Retain the main subject, pose, and overall concept",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@oggii_0",
    "tags": [
      "Illustration & Art",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-299",
    "locale": "zh",
    "title": "极简留白涂鸦手绘草图",
    "category": "插画与艺术",
    "description": "[中文] 以涂鸦速写风表现【主题/主体】，整体呈现快速勾勒、自由变形、即兴手绘与草稿式的视觉效果。线条随手、夸张、可粗细不一，略显凌乱但具有节奏和表现力，强调概括、夸张、趣味和随性，而不是严谨写实或精细刻画。 颜色采用粗糙、干刷感明显的块面表现，可保留不均匀的涂抹痕迹、刷痕、飞白",
    "promptPreview": "[中文] 以涂鸦速写风表现【主题/主体】，整体呈现快速勾勒、自由变形、即兴手绘与草稿式的视觉效果。线条随手、夸张、可粗细不一，略显凌乱但具有节奏和表现力，强调概括、夸张、趣味和随性，而不是严谨写实或精细刻画。 颜色采用粗糙、干刷感明显的块面表现，可保留不均匀的涂抹痕迹、刷痕、飞白与覆盖感，色彩根据【主题/主体】自动适配",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@VoxcatAI",
    "tags": [
      "Illustration & Art",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-452",
    "locale": "en",
    "title": "极简童话手绘儿童插画",
    "category": "插画与艺术",
    "description": "Transform the photo into a delicate minimalist hand-drawn children’s illustration with a soft whimsical fairy-tale aesthetic. Use simple elo",
    "promptPreview": "Transform the photo into a delicate minimalist hand-drawn children’s illustration with a soft whimsical fairy-tale aesthetic. Use simple elongated shapes, thin ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@MissDelulu9",
    "tags": [
      "Illustration & Art",
      "Realistic",
      "Illustration",
      "Character",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-118",
    "locale": "en",
    "title": "漫画分镜叙事设计",
    "category": "插画与艺术",
    "description": "A high-contrast, black-and-white illustration of an elderly man in a sharp suit, drawing a katana. The man has slicked-back white hair, deep",
    "promptPreview": "A high-contrast, black-and-white illustration of an elderly man in a sharp suit, drawing a katana. The man has slicked-back white hair, deep wrinkles, and an in",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@nicdunz",
    "tags": [
      "Illustration & Art",
      "UI",
      "Illustration",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-60",
    "locale": "en",
    "title": "漫画分镜叙事设计",
    "category": "插画与艺术",
    "description": "{ \"type\": \"5-panel collage\", \"layout\": \"grid with 3 top panels and 2 bottom panels\", \"panels\": [ { \"position\": \"top-left\", \"subject\": \"analo",
    "promptPreview": "{ \"type\": \"5-panel collage\", \"layout\": \"grid with 3 top panels and 2 bottom panels\", \"panels\": [ { \"position\": \"top-left\", \"subject\": \"analog clock\", \"details\":",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Gc_qube",
    "tags": [
      "Illustration & Art",
      "UI",
      "Realistic",
      "Illustration",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-125",
    "locale": "en",
    "title": "电商商品展示设计",
    "category": "插画与艺术",
    "description": "{ \"type\": \"anime production layout sheet\", \"style\": \"traditional colored pencil genga, key animation drawing\", \"subject\": { \"character\": \"{a",
    "promptPreview": "{ \"type\": \"anime production layout sheet\", \"style\": \"traditional colored pencil genga, key animation drawing\", \"subject\": { \"character\": \"{argument name=\\\"chara",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Gc_qube",
    "tags": [
      "Illustration & Art",
      "Product",
      "Character",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-231",
    "locale": "zh",
    "title": "疾风起狂草艺术字体设计",
    "category": "插画与艺术",
    "description": "[中文] 创意艺术字体“纵有疾风起”，秀丽笔手写风格，整体文字横版排列，具有强烈视觉冲击力； 深度融合手写书法笔意，笔触带毛笔书写的粗犷洒脱，如挥毫泼墨的肆意劲道； 起收笔的飞白，顿挫，尽显促销的火爆张力，文字的形态打破规整，笔画的粗细变化； dutch angle，营造出动感冲",
    "promptPreview": "[中文] 创意艺术字体“纵有疾风起”，秀丽笔手写风格，整体文字横版排列，具有强烈视觉冲击力； 深度融合手写书法笔意，笔触带毛笔书写的粗犷洒脱，如挥毫泼墨的肆意劲道； 起收笔的飞白，顿挫，尽显促销的火爆张力，文字的形态打破规整，笔画的粗细变化； dutch angle，营造出动感冲刺的气势，字形呈奔放之势； 重心上扬如蓄",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "OpenNana",
    "tags": [
      "Illustration & Art",
      "Poster",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-504",
    "locale": "en",
    "title": "粗糙涂鸦人像改图",
    "category": "插画与艺术",
    "description": "Turn this photo into a chaotic funny doodle illustration, intentionally messy and low-skill, as if drawn quickly with a cheap marker, crayon",
    "promptPreview": "Turn this photo into a chaotic funny doodle illustration, intentionally messy and low-skill, as if drawn quickly with a cheap marker, crayon, or worn-out felt p",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Shorelyn_",
    "tags": [
      "Illustration & Art",
      "UI",
      "Realistic",
      "Illustration",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-129",
    "locale": "en",
    "title": "绘画艺术风格图",
    "category": "插画与艺术",
    "description": "A watercolor illustration of a children's picture book cover. The main subject is a {argument name=\"character appearance\" default=\"cute furr",
    "promptPreview": "A watercolor illustration of a children's picture book cover. The main subject is a {argument name=\"character appearance\" default=\"cute furry kemonomimi girl wi",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@TlanoVRC",
    "tags": [
      "Illustration & Art",
      "Poster",
      "Illustration",
      "Character",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-40",
    "locale": "en",
    "title": "综合应用场景图",
    "category": "插画与艺术",
    "description": "Create {argument name=\"quantity\" default=\"24\"} LINE stickers of {argument name=\"animals\" default=\"animals\"} in a quirky hand-drawn style. Ta",
    "promptPreview": "Create {argument name=\"quantity\" default=\"24\"} LINE stickers of {argument name=\"animals\" default=\"animals\"} in a quirky hand-drawn style. Target {argument name=",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@midori_tatsuta",
    "tags": [
      "Illustration & Art",
      "UI",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-442",
    "locale": "en",
    "title": "舒适发廊插画",
    "category": "插画与艺术",
    "description": "A vibrant whimsical digital illustration of a cozy indie hair salon, featuring a young girl with long brown hair getting her hair styled by ",
    "promptPreview": "A vibrant whimsical digital illustration of a cozy indie hair salon, featuring a young girl with long brown hair getting her hair styled by a fashionable hairst",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Sairah_0",
    "tags": [
      "Illustration & Art",
      "Illustration",
      "Character",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-471",
    "locale": "en",
    "title": "花簪和服动漫肖像",
    "category": "插画与艺术",
    "description": "Ultra-detailed anime-style portrait of a young girl with large expressive eyes, soft blush cheeks, and delicate facial features. She wears a",
    "promptPreview": "Ultra-detailed anime-style portrait of a young girl with large expressive eyes, soft blush cheeks, and delicate facial features. She wears a vibrant floral kimo",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Mind_Boticni",
    "tags": [
      "Illustration & Art",
      "Illustration",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-233",
    "locale": "zh",
    "title": "蒙娜丽莎畅饮可乐的趣味油画",
    "category": "插画与艺术",
    "description": "[中文] 生成一张蒙娜丽莎喝可乐的油画。 [English] Generate an oil painting of Mona Lisa drinking cola.",
    "promptPreview": "[中文] 生成一张蒙娜丽莎喝可乐的油画。 [English] Generate an oil painting of Mona Lisa drinking cola.",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "Illustration & Art",
      "Illustration",
      "Tech",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-281",
    "locale": "zh",
    "title": "赛博朋克科幻曼荼罗",
    "category": "插画与艺术",
    "description": "[中文] でChatGPTで画像を作成してもらって、今日また作成してもらったらGPT image 2かもしれず、出来が変わったように見えるのでメモ 左の水色と黄色のが先週 右の紫のが今日 右のは透明感とか解像度、緻密さが違うような気がする… プロンプト 曼荼羅の近未来SF版を描い",
    "promptPreview": "[中文] でChatGPTで画像を作成してもらって、今日また作成してもらったらGPT image 2かもしれず、出来が変わったように見えるのでメモ 左の水色と黄色のが先週 右の紫のが今日 右のは透明感とか解像度、緻密さが違うような気がする… プロンプト 曼荼羅の近未来SF版を描いて [English] I had Ch",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@4WEB1",
    "tags": [
      "Illustration & Art",
      "Illustration",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-356",
    "locale": "en",
    "title": "过度思考超现实街头 Campaign",
    "category": "插画与艺术",
    "description": "Ultra-realistic conceptual portrait of a young woman with long wavy hair and soft defined features, wearing rose-tinted rectangular sunglass",
    "promptPreview": "Ultra-realistic conceptual portrait of a young woman with long wavy hair and soft defined features, wearing rose-tinted rectangular sunglasses, an oversized ivo",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@AIwithAliya",
    "tags": [
      "Illustration & Art",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-498",
    "locale": "en",
    "title": "铅笔画背景 3D 分身",
    "category": "插画与艺术",
    "description": "Create a handrawn pencil illustration of [image] yawning on paper, as background. Add a 3D Pixar style render of [foto] standing casually in",
    "promptPreview": "Create a handrawn pencil illustration of [image] yawning on paper, as background. Add a 3D Pixar style render of [foto] standing casually infront of the giant h",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@AIwithkhan",
    "tags": [
      "Illustration & Art",
      "Illustration",
      "3D",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-430",
    "locale": "en",
    "title": "铅笔素描时尚编辑插画",
    "category": "插画与艺术",
    "description": "A high-detail digital illustration of a stylish woman sitting gracefully on a stone ledge, posing with one hand near her chin and legs cross",
    "promptPreview": "A high-detail digital illustration of a stylish woman sitting gracefully on a stone ledge, posing with one hand near her chin and legs crossed. She is wearing v",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@harboriis",
    "tags": [
      "Illustration & Art",
      "Illustration",
      "3D",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-433",
    "locale": "en",
    "title": "韩国城市水彩旅行插画",
    "category": "插画与艺术",
    "description": "Dreamy watercolor travel illustration of a peaceful Korean city street, hand-painted urban sketchbook style, delicate ink linework mixed wit",
    "promptPreview": "Dreamy watercolor travel illustration of a peaceful Korean city street, hand-painted urban sketchbook style, delicate ink linework mixed with soft watercolor wa",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Taaruk_",
    "tags": [
      "Illustration & Art",
      "UI",
      "Poster",
      "Illustration",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-246",
    "locale": "zh",
    "title": "黑白线稿勾勒的上海风情",
    "category": "插画与艺术",
    "description": "[中文] 设计一张黑色线稿风格的上海明信片 [English] Design a Shanghai postcard in black line art style.",
    "promptPreview": "[中文] 设计一张黑色线稿风格的上海明信片 [English] Design a Shanghai postcard in black line art style.",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@akokoi1",
    "tags": [
      "Illustration & Art",
      "Illustration",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-427",
    "locale": "en",
    "title": "9-frame 时尚人像拼贴",
    "category": "摄影与写实",
    "description": "Edit this photo and don't change the face, portrait 9:16. A 9-frame fashion portrait collage featuring a stylish young woman with playful ex",
    "promptPreview": "Edit this photo and don't change the face, portrait 9:16. A 9-frame fashion portrait collage featuring a stylish young woman with playful expressions and sporty",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@saniaspeaks_",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Character",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-383",
    "locale": "en",
    "title": "AI 日常生活 iPhone 抓拍",
    "category": "摄影与写实",
    "description": "I want to see what you really look like. Draw a snapshot of your everyday life as if it were accidentally taken on an iPhone. Make it feel l",
    "promptPreview": "I want to see what you really look like. Draw a snapshot of your everyday life as if it were accidentally taken on an iPhone. Make it feel like a very ordinary,",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Ciri_ai",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-408",
    "locale": "en",
    "title": "Cozy Academia 学习手记",
    "category": "摄影与写实",
    "description": "Dreamy cinematic study aesthetic, young Asian girl with long dark hair studying outdoors at a wooden table during golden hour, cozy oversize",
    "promptPreview": "Dreamy cinematic study aesthetic, young Asian girl with long dark hair studying outdoors at a wooden table during golden hour, cozy oversized green sweater and ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Sairah_0",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Product",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-428",
    "locale": "en",
    "title": "F1 直播转播围场截图",
    "category": "摄影与写实",
    "description": "Ultra-realistic F1 live TV broadcast screenshot. A strikingly beautiful, 25-year-old young woman with bright light blue hair and captivating",
    "promptPreview": "Ultra-realistic F1 live TV broadcast screenshot. A strikingly beautiful, 25-year-old young woman with bright light blue hair and captivating bright cerulean blu",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@bigwonbots",
    "tags": [
      "Photography & Realism",
      "UI",
      "Realistic",
      "Brand",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-421",
    "locale": "en",
    "title": "iPhone 屏幕遮脸创意人像",
    "category": "摄影与写实",
    "description": "Ultra-realistic creative portrait taken with an iPhone, identity accurately preserved from the reference image. A woman stands inside a stor",
    "promptPreview": "Ultra-realistic creative portrait taken with an iPhone, identity accurately preserved from the reference image. A woman stands inside a store, facing a glass di",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Ciri_ai",
    "tags": [
      "Photography & Realism",
      "Poster",
      "Realistic",
      "Brand",
      "Tech",
      "Fashion"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-393",
    "locale": "en",
    "title": "Y2K 金色时刻人像",
    "category": "摄影与写实",
    "description": "candid portrait of a beautiful young blonde woman, 21 years old, glowing sun-kissed skin, wearing a baby pink velour tracksuit and butterfly",
    "promptPreview": "candid portrait of a beautiful young blonde woman, 21 years old, glowing sun-kissed skin, wearing a baby pink velour tracksuit and butterfly clips in her hair, ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@SiliconBarbie_",
    "tags": [
      "Photography & Realism",
      "UI",
      "Realistic",
      "Tech",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-491",
    "locale": "en",
    "title": "Y2K 高楼浴室镜面自拍",
    "category": "摄影与写实",
    "description": "A moody late-night mirror selfie captured in a luxury high-rise bathroom overlooking a glowing city skyline. A young woman with long, slight",
    "promptPreview": "A moody late-night mirror selfie captured in a luxury high-rise bathroom overlooking a glowing city skyline. A young woman with long, slightly messy black hair ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@jzaib4269",
    "tags": [
      "Photography & Realism",
      "UI",
      "Realistic",
      "Product",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-472",
    "locale": "en",
    "title": "上海地铁站台晨光",
    "category": "摄影与写实",
    "description": "A candid photograph of a young woman standing on a Shanghai metro platform during the summer morning commute, authentic daily life photograp",
    "promptPreview": "A candid photograph of a young woman standing on a Shanghai metro platform during the summer morning commute, authentic daily life photography, natural candid m",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@ToroJushiAi",
    "tags": [
      "Photography & Realism",
      "Poster",
      "Realistic",
      "Illustration",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-45",
    "locale": "en",
    "title": "人像写实摄影图",
    "category": "摄影与写实",
    "description": "A striking black and white close-up portrait of a {argument name=\"subject description\" default=\"handsome young Asian man\"} with {argument na",
    "promptPreview": "A striking black and white close-up portrait of a {argument name=\"subject description\" default=\"handsome young Asian man\"} with {argument name=\"hair style\" defa",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@AoYe999",
    "tags": [
      "Photography & Realism",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-35",
    "locale": "en",
    "title": "人像写实摄影图",
    "category": "摄影与写实",
    "description": "A {argument name=\"photography style\" default=\"photorealistic portrait with shallow depth of field and soft bokeh\"} of a {argument name=\"subj",
    "promptPreview": "A {argument name=\"photography style\" default=\"photorealistic portrait with shallow depth of field and soft bokeh\"} of a {argument name=\"subject\" default=\"young ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@kazmaendo",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Tech",
      "Fashion",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-154",
    "locale": "en",
    "title": "写实摄影风格创作",
    "category": "摄影与写实",
    "description": "A photorealistic, high-resolution commercial photograph of a {argument name=\"car model and color\" default=\"bright blue Alpine A110 R sports ",
    "promptPreview": "A photorealistic, high-resolution commercial photograph of a {argument name=\"car model and color\" default=\"bright blue Alpine A110 R sports car\"} parked in the ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@AlwaveNazca",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-56",
    "locale": "en",
    "title": "写实摄影风格创作",
    "category": "摄影与写实",
    "description": "A candid, realistic photograph of a young {argument name=\"subject aesthetic\" default=\"goth\"} woman with pale skin, long straight black hair ",
    "promptPreview": "A candid, realistic photograph of a young {argument name=\"subject aesthetic\" default=\"goth\"} woman with pale skin, long straight black hair with bangs, heavy bl",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@danieldmai",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-28",
    "locale": "en",
    "title": "写实摄影风格创作",
    "category": "摄影与写实",
    "description": "{ \"type\": \"2x2 portrait grid\", \"subject\": \"{argument name=\\\"subject description\\\" default=\\\"young adult East Asian male with short black hai",
    "promptPreview": "{ \"type\": \"2x2 portrait grid\", \"subject\": \"{argument name=\\\"subject description\\\" default=\\\"young adult East Asian male with short black hair and a slight smile",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@frankfu1688",
    "tags": [
      "Photography & Realism",
      "UI",
      "Realistic",
      "Brand",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-490",
    "locale": "en",
    "title": "双重曝光时尚肖像",
    "category": "摄影与写实",
    "description": "A cinematic double-exposure fashion portrait featuring a woman with the uploaded face used 100% as reference, wearing stylish metallic aviat",
    "promptPreview": "A cinematic double-exposure fashion portrait featuring a woman with the uploaded face used 100% as reference, wearing stylish metallic aviator sunglasses. The c",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Anaya_Ai12",
    "tags": [
      "Photography & Realism",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-366",
    "locale": "en",
    "title": "咖啡馆写实照片与 2D 涂鸦叠加",
    "category": "摄影与写实",
    "description": "A trendy young woman sitting at an outdoor café table, holding a hot coffee cup near her lips. She has short curly black hair and wears over",
    "promptPreview": "A trendy young woman sitting at an outdoor café table, holding a hot coffee cup near her lips. She has short curly black hair and wears oversized tinted sunglas",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Jawad_Rahman_",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Illustration",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-143",
    "locale": "en",
    "title": "品牌徽标设计图",
    "category": "摄影与写实",
    "description": "A photorealistic amateur photograph of a custom building block set resting on a light wood grain table in a living room. In the background s",
    "promptPreview": "A photorealistic amateur photograph of a custom building block set resting on a light wood grain table in a living room. In the background stands a large produc",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Gc_qube",
    "tags": [
      "Photography & Realism",
      "UI",
      "Realistic",
      "Product",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-36",
    "locale": "en",
    "title": "品牌徽标设计图",
    "category": "摄影与写实",
    "description": "A photorealistic selfie of a young man with short wavy dark hair and light stubble on an indoor basketball court. He wears a black athletic ",
    "promptPreview": "A photorealistic selfie of a young man with short wavy dark hair and light stubble on an indoor basketball court. He wears a black athletic t-shirt with a white",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@mirochill",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Brand",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-501",
    "locale": "en",
    "title": "夏日牵手回眸电影肖像",
    "category": "摄影与写实",
    "description": "Cinematic portrait photography, ultra-photorealistic, 2160x3840 vertical composition, 50mm or 85mm portrait lens rendering, shallow depth of",
    "promptPreview": "Cinematic portrait photography, ultra-photorealistic, 2160x3840 vertical composition, 50mm or 85mm portrait lens rendering, shallow depth of field, clean transl",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@iamaiistudio",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Character",
      "Classical",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-277",
    "locale": "zh",
    "title": "奢华魅力黑人女性海滨摄影",
    "category": "摄影与写实",
    "description": "[中文] 奢华魅力美容肖像：, 美丽的黑人女性, 青春活力, 奶油香草色, 丝绸柔顺发, 红木色, 微妙的自信, 有质感的面料, 蓝宝石色, 极简珠宝, 海滨微风, 镜头光晕效果, 怀旧的, 电影镜头, 对称构图, 柔焦, 高级时尚摄影, 单色的, 水光质感, 神秘张力, 分层元",
    "promptPreview": "[中文] 奢华魅力美容肖像：, 美丽的黑人女性, 青春活力, 奶油香草色, 丝绸柔顺发, 红木色, 微妙的自信, 有质感的面料, 蓝宝石色, 极简珠宝, 海滨微风, 镜头光晕效果, 怀旧的, 电影镜头, 对称构图, 柔焦, 高级时尚摄影, 单色的, 水光质感, 神秘张力, 分层元素 [English] Luxury G",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@patrickassale",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Tech",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-202",
    "locale": "zh",
    "title": "宅男必看绝美二次元少女",
    "category": "摄影与写实",
    "description": "[中文] 生成高质量美女（宅男必备） [English] Generate high-quality beautiful girl (otaku must-have)",
    "promptPreview": "[中文] 生成高质量美女（宅男必备） [English] Generate high-quality beautiful girl (otaku must-have)",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@joshesye",
    "tags": [
      "Photography & Realism",
      "Photography",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-414",
    "locale": "en",
    "title": "室内晨间写实摄影",
    "category": "摄影与写实",
    "description": "A close-medium shot of a young Japanese woman in her bedroom on an ordinary morning, captured in authentic daily life photography as a natur",
    "promptPreview": "A close-medium shot of a young Japanese woman in her bedroom on an ordinary morning, captured in authentic daily life photography as a natural candid moment. Sh",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@ToroJushiAi",
    "tags": [
      "Photography & Realism",
      "UI",
      "Realistic",
      "Illustration",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-53",
    "locale": "en",
    "title": "室内空间渲染图",
    "category": "摄影与写实",
    "description": "A vintage, late 90s amateur flash photograph of a young man repairing an arcade machine. He is kneeling on a dark, patterned arcade carpet, ",
    "promptPreview": "A vintage, late 90s amateur flash photograph of a young man repairing an arcade machine. He is kneeling on a dark, patterned arcade carpet, looking back over hi",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@nicdunz",
    "tags": [
      "Photography & Realism",
      "UI",
      "Realistic",
      "Brand",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-488",
    "locale": "en",
    "title": "屋顶球场日落人像",
    "category": "摄影与写实",
    "description": "Ultra-realistic, high-quality portrait of a stylish young South Asian woman standing gracefully on a colorful urban rooftop basketball court",
    "promptPreview": "Ultra-realistic, high-quality portrait of a stylish young South Asian woman standing gracefully on a colorful urban rooftop basketball court at sunset. She is s",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@HaniaAi12",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Tech",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-412",
    "locale": "en",
    "title": "彩色按钮时尚 Campaign",
    "category": "摄影与写实",
    "description": "Use reference image as style guide. Surreal minimalist fashion editorial photography with vibrant oversized colorful buttons as main visual ",
    "promptPreview": "Use reference image as style guide. Surreal minimalist fashion editorial photography with vibrant oversized colorful buttons as main visual theme (cyan, orange,",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Mind_Boticni",
    "tags": [
      "Photography & Realism",
      "UI",
      "Realistic",
      "Character",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-426",
    "locale": "en",
    "title": "日韩咖啡馆情侣写真",
    "category": "摄影与写实",
    "description": "Ultra-realistic cozy Japanese-Korean café photography featuring a cute young [Japanese/Korean] couple sitting together naturally in a trendy",
    "promptPreview": "Ultra-realistic cozy Japanese-Korean café photography featuring a cute young [Japanese/Korean] couple sitting together naturally in a trendy aesthetic café. The",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@sha_zdiii",
    "tags": [
      "Photography & Realism",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-382",
    "locale": "en",
    "title": "春日花田三联竖版写真拼贴",
    "category": "摄影与写实",
    "description": "A high-quality 3-panel vertical photo collage of a stunning uploaded woman with soft, voluminous wavy hair glowing in golden sunlight. She i",
    "promptPreview": "A high-quality 3-panel vertical photo collage of a stunning uploaded woman with soft, voluminous wavy hair glowing in golden sunlight. She is styled in a cream-",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@frametheory058",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Character",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-311",
    "locale": "zh",
    "title": "晨曦薰衣草田梦幻少女三联画",
    "category": "摄影与写实",
    "description": "[中文] 日出时分薰衣草田中女子的水平三联画。 上部：半身像，闭着眼睛，淡紫色连衣裙，一只手放在头发里，模糊的薰衣草前景。 中部：特写镜头，看着镜头，蓬乱的头发，薄纱围巾，脸上的阳光。 下部：四分之三镜头，手持薰衣草花束，飘逸的裙子，柔和的粉彩天空，温暖的梦幻色调。 [Engli",
    "promptPreview": "[中文] 日出时分薰衣草田中女子的水平三联画。 上部：半身像，闭着眼睛，淡紫色连衣裙，一只手放在头发里，模糊的薰衣草前景。 中部：特写镜头，看着镜头，蓬乱的头发，薄纱围巾，脸上的阳光。 下部：四分之三镜头，手持薰衣草花束，飘逸的裙子，柔和的粉彩天空，温暖的梦幻色调。 [English] Horizontal tript",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@Naiknelofar788",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-499",
    "locale": "en",
    "title": "极简精品店全身时尚写真",
    "category": "摄影与写实",
    "description": "A full-body editorial fashion photograph of a beautiful young woman with the same appearance as the reference, long glossy dark hair, soft b",
    "promptPreview": "A full-body editorial fashion photograph of a beautiful young woman with the same appearance as the reference, long glossy dark hair, soft bangs, fair skin, and",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@john_my07",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Character",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-500",
    "locale": "en",
    "title": "梦幻花冠仙境肖像",
    "category": "摄影与写实",
    "description": "Ultra-realistic ethereal fantasy portrait of a breathtaking young woman with delicate porcelain skin, soft grey-blue eyes, and natural rosy ",
    "promptPreview": "Ultra-realistic ethereal fantasy portrait of a breathtaking young woman with delicate porcelain skin, soft grey-blue eyes, and natural rosy lips. She gazes gent",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@HaniaAi12",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-273",
    "locale": "zh",
    "title": "橙红渐变中的孤独剪影",
    "category": "摄影与写实",
    "description": "[中文] 生成一张电影级极简肖像，一个孤独的男人站在强烈的橙色到红色渐变环境中，强烈的剪影光，深邃的阴影对比，反光的光滑地面，对称构图，极简 [English] Generate a cinematic minimal portrait of a solitary man sta",
    "promptPreview": "[中文] 生成一张电影级极简肖像，一个孤独的男人站在强烈的橙色到红色渐变环境中，强烈的剪影光，深邃的阴影对比，反光的光滑地面，对称构图，极简 [English] Generate a cinematic minimal portrait of a solitary man standing in an intense ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@iam_miharbi",
    "tags": [
      "Photography & Realism",
      "Photography",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-467",
    "locale": "zh",
    "title": "泳装杂志九宫格广告页",
    "category": "摄影与写实",
    "description": "泳装时尚杂志广告页面，日本成熟模特，S型曲线。变换姿势和风格，九宫格展示，保持人物面部一致性",
    "promptPreview": "泳装时尚杂志广告页面，日本成熟模特，S型曲线。变换姿势和风格，九宫格展示，保持人物面部一致性",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@Adam38363368936",
    "tags": [
      "Photography & Realism",
      "Character",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-165",
    "locale": "zh",
    "title": "清冷佳人夜市烧烤三刀流",
    "category": "摄影与写实",
    "description": "[中文] 一个有着清冷孤傲气质的绝美佳人，精致的面部特征，一张冷酷且精致的高级时装面容，长发，以及优雅苗条的身材；烧烤“三刀流”姿势：嘴里叼着一根烧烤串，每只手各拿一根烧烤串交叉以模仿索隆的三刀流；街头夜景氛围，温暖黄色的夜市灯光，模糊的背景，胶片般的质感，柔焦光晕，电影般的叙事",
    "promptPreview": "[中文] 一个有着清冷孤傲气质的绝美佳人，精致的面部特征，一张冷酷且精致的高级时装面容，长发，以及优雅苗条的身材；烧烤“三刀流”姿势：嘴里叼着一根烧烤串，每只手各拿一根烧烤串交叉以模仿索隆的三刀流；街头夜景氛围，温暖黄色的夜市灯光，模糊的背景，胶片般的质感，柔焦光晕，电影般的叙事感，时髦高端网红风格的时尚拍摄，清晰发光",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@BubbleBrain",
    "tags": [
      "Photography & Realism",
      "Character",
      "Tech",
      "Fashion",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-284",
    "locale": "zh",
    "title": "温馨卧室里的少女自拍",
    "category": "摄影与写实",
    "description": "[中文] 一个惊艳的18岁中国女孩，拥有年轻纯真的脸庞和逼真的皮肤纹理，坐在她卧室里一张舒适且略显凌乱的床上。她正用智能手机拍镜子自拍，捕捉一个自然且亲密的瞬间。穿着随意的灰色居家服和整洁的白色船袜。柔和的自然光（黄金时刻）从侧面窗户照进来，营造出一种温暖、富有情绪感和电影般的氛",
    "promptPreview": "[中文] 一个惊艳的18岁中国女孩，拥有年轻纯真的脸庞和逼真的皮肤纹理，坐在她卧室里一张舒适且略显凌乱的床上。她正用智能手机拍镜子自拍，捕捉一个自然且亲密的瞬间。穿着随意的灰色居家服和整洁的白色船袜。柔和的自然光（黄金时刻）从侧面窗户照进来，营造出一种温暖、富有情绪感和电影般的氛围。35毫米镜头，对镜子中的主体保持锐利",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@Shinning1010",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-24",
    "locale": "en",
    "title": "漫画分镜叙事设计",
    "category": "摄影与写实",
    "description": "Genshin Impact {argument name=\"character\" default=\"Raiden Shogun\"} cosplay selfies at the {argument name=\"event\" default=\"Shanghai Comic Con",
    "promptPreview": "Genshin Impact {argument name=\"character\" default=\"Raiden Shogun\"} cosplay selfies at the {argument name=\"event\" default=\"Shanghai Comic Con\"}",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@wewe50770964683",
    "tags": [
      "Photography & Realism",
      "Character",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-450",
    "locale": "en",
    "title": "烛光侧室写实摄影",
    "category": "摄影与写实",
    "description": "A medium-shot authentic daily life photograph, natural candid moment, of an East Asian young woman standing in a small side room lit only by",
    "promptPreview": "A medium-shot authentic daily life photograph, natural candid moment, of an East Asian young woman standing in a small side room lit only by a single candle on ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@ToroJushiAi",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Illustration",
      "3D",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-29",
    "locale": "en",
    "title": "电影感叙事场景图",
    "category": "摄影与写实",
    "description": "Using REFERENCE_0, transform the subject's appearance to a {argument name=\"style\" default=\"trad goth\"} aesthetic while preserving the exact ",
    "promptPreview": "Using REFERENCE_0, transform the subject's appearance to a {argument name=\"style\" default=\"trad goth\"} aesthetic while preserving the exact pose, clothing struc",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@danieldmai",
    "tags": [
      "Photography & Realism",
      "Character",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-221",
    "locale": "zh",
    "title": "窗边日系胶片女孩",
    "category": "摄影与写实",
    "description": "[中文] 模拟35毫米胶片摄影，柔和轻盈的日系美学，温柔漫射的自然窗户光，轻微过曝，柔和色调，低对比度，柔和的高光，靠近窗户配有白色窗帘的极简室内环境，干净的浅色墙壁，自然构图，平视视角，略微紧凑的全身取景（大腿中部到头部），年轻东亚女性，自然极简妆容，柔和真实的皮肤纹理，长长的",
    "promptPreview": "[中文] 模拟35毫米胶片摄影，柔和轻盈的日系美学，温柔漫射的自然窗户光，轻微过曝，柔和色调，低对比度，柔和的高光，靠近窗户配有白色窗帘的极简室内环境，干净的浅色墙壁，自然构图，平视视角，略微紧凑的全身取景（大腿中部到头部），年轻东亚女性，自然极简妆容，柔和真实的皮肤纹理，长长的微乱黑发，超大号白色纽扣衬衫，浅色休闲短",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@BubbleBrain",
    "tags": [
      "Photography & Realism",
      "UI",
      "Realistic",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-420",
    "locale": "en",
    "title": "红跑道低角度夏日人像",
    "category": "摄影与写实",
    "description": "Use the uploaded portrait photo as the appearance reference for the person. Create a bright photorealistic outdoor portrait of a young woman",
    "promptPreview": "Use the uploaded portrait photo as the appearance reference for the person. Create a bright photorealistic outdoor portrait of a young woman lying on a red runn",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Shinning1010",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Brand",
      "Character",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-482",
    "locale": "en",
    "title": "自我凝视超现实 Campaign",
    "category": "摄影与写实",
    "description": "Ultra-realistic conceptual portrait of a young man with curly hair and light stubble, wearing yellow-tinted rectangular sunglasses, a beige ",
    "promptPreview": "Ultra-realistic conceptual portrait of a young man with curly hair and light stubble, wearing yellow-tinted rectangular sunglasses, a beige minimal t-shirt, blu",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Shorelyn_",
    "tags": [
      "Photography & Realism",
      "Poster",
      "Realistic",
      "Brand",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-198",
    "locale": "zh",
    "title": "苍白陶瓷娃娃沙滩仰视",
    "category": "摄影与写实",
    "description": "[中文] { \"相机参数\": { \"设备类型\": \"iPhone 15 Pro 前置自拍\", \"镜头\": \"24mm\", \"构图\": \"高角度 POV（第一人称视角）\", \"后期处理\": \"计算摄影风格，清晰的数字读出，深景深\" }, \"主体描述\": { \"特征\": \"陶瓷娃娃审",
    "promptPreview": "[中文] { \"相机参数\": { \"设备类型\": \"iPhone 15 Pro 前置自拍\", \"镜头\": \"24mm\", \"构图\": \"高角度 POV（第一人称视角）\", \"后期处理\": \"计算摄影风格，清晰的数字读出，深景深\" }, \"主体描述\": { \"特征\": \"陶瓷娃娃审美，无瑕的苍白皮肤，巨大的冰蓝色眼睛，小",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@IamEmily2050",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-262",
    "locale": "zh",
    "title": "苹果园远观库克发布新机",
    "category": "摄影与写实",
    "description": "[中文] 在Apple Park iPhone 20主题演讲期间拍摄的业余iPhone照片，蒂姆·库克在舞台上演讲。从远处的观众人群中拍摄 [English] Amateur iPhone photo at Apple Park during the iPhone 20 keyn",
    "promptPreview": "[中文] 在Apple Park iPhone 20主题演讲期间拍摄的业余iPhone照片，蒂姆·库克在舞台上演讲。从远处的观众人群中拍摄 [English] Amateur iPhone photo at Apple Park during the iPhone 20 keynote, Tim Cook presen",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@austinit",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-322",
    "locale": "zh",
    "title": "街头炫瓶男模",
    "category": "摄影与写实",
    "description": "[中文] 专业照片，一位男士，30岁的俄罗斯模特（参考图像），正对着镜头，向相机倾斜，从下往上拍摄，使用广角镜头。男士倾斜着身体，近距离将一瓶饮料展示给镜头，一只手拿着瓶子，紧贴在镜头前。瓶子的标签和方向保持笔直，以便标签清晰可读。他穿着白色运动鞋，一只脚在镜头前方。男士站在街道",
    "promptPreview": "[中文] 专业照片，一位男士，30岁的俄罗斯模特（参考图像），正对着镜头，向相机倾斜，从下往上拍摄，使用广角镜头。男士倾斜着身体，近距离将一瓶饮料展示给镜头，一只手拿着瓶子，紧贴在镜头前。瓶子的标签和方向保持笔直，以便标签清晰可读。他穿着白色运动鞋，一只脚在镜头前方。男士站在街道上，湿漉漉的沥青和飞溅的水花从下方拍出。",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@ecommartinez",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-195",
    "locale": "zh",
    "title": "超写实与水墨的梦幻融合",
    "category": "摄影与写实",
    "description": "[中文] 一张动态的混合媒体摄影作品，将超写实肖像与传统的中国水墨插画相融合。 中心人物是一位具有柔和短波波头短发造型的照片般逼真的年轻亚洲女性。她的妆容自然且极简，表情平静而温柔。她背对相机站立，姿态呈现出优雅的S型曲线，营造出优美流畅的剪影。她微微转动上半身，越过肩膀回眸，带",
    "promptPreview": "[中文] 一张动态的混合媒体摄影作品，将超写实肖像与传统的中国水墨插画相融合。 中心人物是一位具有柔和短波波头短发造型的照片般逼真的年轻亚洲女性。她的妆容自然且极简，表情平静而温柔。她背对相机站立，姿态呈现出优雅的S型曲线，营造出优美流畅的剪影。她微微转动上半身，越过肩膀回眸，带着一种安静、内省的情绪。 她穿着一件简约",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@johnAGI168",
    "tags": [
      "Photography & Realism",
      "UI",
      "Realistic",
      "Illustration",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-199",
    "locale": "zh",
    "title": "超写实海滩高角度手机自拍",
    "category": "摄影与写实",
    "description": "[中文] { 超写实iPhone 15 Pro前置摄像头自拍，一位成年女性在明亮的沙滩上， 从举臂高角度自拍视角拍摄。手机略微举在脸部上方， 营造出自然的前置摄像头几何形态，带有轻微的等效24mm广角畸变， 写实的面部比例， 以及智能手机的深景深。她向上抬起下巴，一只手遮挡刺眼的",
    "promptPreview": "[中文] { 超写实iPhone 15 Pro前置摄像头自拍，一位成年女性在明亮的沙滩上， 从举臂高角度自拍视角拍摄。手机略微举在脸部上方， 营造出自然的前置摄像头几何形态，带有轻微的等效24mm广角畸变， 写实的面部比例， 以及智能手机的深景深。她向上抬起下巴，一只手遮挡刺眼的阳光，同时直视手机镜头。她的表情中性， ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@IamEmily2050",
    "tags": [
      "Photography & Realism",
      "UI",
      "Realistic",
      "3D",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-349",
    "locale": "en",
    "title": "运动时尚三联 Campaign",
    "category": "摄影与写实",
    "description": "Cinematic sports fashion collage, 3-panel layout, top panel large hero shot of a female tennis athlete sitting confidently on an oversized t",
    "promptPreview": "Cinematic sports fashion collage, 3-panel layout, top panel large hero shot of a female tennis athlete sitting confidently on an oversized tilted tennis racket,",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@AIwithkhan",
    "tags": [
      "Photography & Realism",
      "Poster",
      "Brand",
      "Character",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-483",
    "locale": "en",
    "title": "都市飞鸟街头肖像",
    "category": "摄影与写实",
    "description": "A stylish cinematic portrait of a confident young woman leaning casually against a textured urban concrete wall, surrounded by vibrant flyin",
    "promptPreview": "A stylish cinematic portrait of a confident young woman leaning casually against a textured urban concrete wall, surrounded by vibrant flying birds including bl",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@HaniaAi12",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-429",
    "locale": "en",
    "title": "韩国便利店粉色 Hoodie 人像",
    "category": "摄影与写实",
    "description": "Ultra-realistic cozy Korean convenience store portrait of a beautiful Korean woman standing in front of glowing refrigerator aisles at night",
    "promptPreview": "Ultra-realistic cozy Korean convenience store portrait of a beautiful Korean woman standing in front of glowing refrigerator aisles at night, wearing an oversiz",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@john_my07",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Product",
      "Tech",
      "Fashion",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-219",
    "locale": "zh",
    "title": "韩系偶像九宫格写真集",
    "category": "摄影与写实",
    "description": "[中文] 9:16 竖版 — 一个 3x3 网格拼贴（九张图片）形成一系列韩国偶像肖像摄影。每一帧都呈现同一位年轻的韩国女性偶像，在所有九张镜头中保持 100% 一致的面部特征、比例、发型和身份。自然、超逼真的皮肤纹理，无修图，无磨皮。干净的偶像风格极简妆容，柔和的光泽，微妙的瑕",
    "promptPreview": "[中文] 9:16 竖版 — 一个 3x3 网格拼贴（九张图片）形成一系列韩国偶像肖像摄影。每一帧都呈现同一位年轻的韩国女性偶像，在所有九张镜头中保持 100% 一致的面部特征、比例、发型和身份。自然、超逼真的皮肤纹理，无修图，无磨皮。干净的偶像风格极简妆容，柔和的光泽，微妙的瑕疵。发型：长发、蓬松的黑发，微乱，在所有",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@BubbleBrain",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Brand",
      "Tech",
      "Fashion",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-466",
    "locale": "en",
    "title": "鱼市追猫 CCD 街拍",
    "category": "摄影与写实",
    "description": "Use the uploaded reference image as the exact identity base for the main subject. Preserve their authentic facial structure, recognizable ap",
    "promptPreview": "Use the uploaded reference image as the exact identity base for the main subject. Preserve their authentic facial structure, recognizable appearance, hairstyle,",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@mehvishs25",
    "tags": [
      "Photography & Realism",
      "UI",
      "Realistic",
      "Brand",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-357",
    "locale": "en",
    "title": "鱼眼镜面复古咖啡馆人像",
    "category": "摄影与写实",
    "description": "A fish-eye lens close-up of [your photo as reference] sipping from a teal/turquoise coffee mug, leaning forward intimately toward camera. Sh",
    "promptPreview": "A fish-eye lens close-up of [your photo as reference] sipping from a teal/turquoise coffee mug, leaning forward intimately toward camera. Shot through or near a",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@harboriis",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Food"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-300",
    "locale": "zh",
    "title": "黑板上的出师表全文",
    "category": "摄影与写实",
    "description": "[中文] 生成图片: 手写在教室黑板上的出师表全文，真实感的粉笔字迹，晴朗白天用iPhone手机实拍 [English] Generate image: The full text of Chu Shi Biao handwritten on a classroom blackb",
    "promptPreview": "[中文] 生成图片: 手写在教室黑板上的出师表全文，真实感的粉笔字迹，晴朗白天用iPhone手机实拍 [English] Generate image: The full text of Chu Shi Biao handwritten on a classroom blackboard, realistic chal",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@rionaifantasy",
    "tags": [
      "Photography & Realism",
      "Realistic",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-425",
    "locale": "en",
    "title": "黑白时尚人像拼贴海报",
    "category": "摄影与写实",
    "description": "{ \"prompt\": \"Vertical poster collage design, beige background, four stacked horizontal rounded rectangles containing black and white cinemat",
    "promptPreview": "{ \"prompt\": \"Vertical poster collage design, beige background, four stacked horizontal rounded rectangles containing black and white cinematic portraits of the ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@XSydneyFan",
    "tags": [
      "Photography & Realism",
      "Poster",
      "Realistic",
      "Character",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-492",
    "locale": "en",
    "title": "黑色高定酒店套房写真",
    "category": "摄影与写实",
    "description": "Create a hyper-realistic luxury fashion portrait inspired by premium fashion-week editorials. IMPORTANT: Generate ONLY ONE SINGLE VERTICAL I",
    "promptPreview": "Create a hyper-realistic luxury fashion portrait inspired by premium fashion-week editorials. IMPORTANT: Generate ONLY ONE SINGLE VERTICAL IMAGE. STYLE: Luxury ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@xRahultripathi",
    "tags": [
      "Photography & Realism",
      "UI",
      "Realistic",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-201",
    "locale": "zh",
    "title": "三甲医院真实门诊处方笺",
    "category": "文档与出版",
    "description": "[中文] 一张三甲医院的门诊处方笺，医生潦草的手写字，包含真实合理的 诊断、药品名、剂量，右下角有医生签名和科室章。 [English] An outpatient prescription sheet from a Grade 3A hospital, doctor's ill",
    "promptPreview": "[中文] 一张三甲医院的门诊处方笺，医生潦草的手写字，包含真实合理的 诊断、药品名、剂量，右下角有医生签名和科室章。 [English] An outpatient prescription sheet from a Grade 3A hospital, doctor's illegible handwriting, ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@msjiaozhu",
    "tags": [
      "Documents & Publishing",
      "Realistic",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-119",
    "locale": "en",
    "title": "主题海报版式设计",
    "category": "文档与出版",
    "description": "{ \"type\": \"anime movie production pitch document\", \"overall_layout\": \"split layout with a large cinematic movie poster on the top half and a",
    "promptPreview": "{ \"type\": \"anime movie production pitch document\", \"overall_layout\": \"split layout with a large cinematic movie poster on the top half and a grid of 5 detailed ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@old_pgmrs_will",
    "tags": [
      "Documents & Publishing",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-303",
    "locale": "zh",
    "title": "人教版三年级语文课本内页",
    "category": "文档与出版",
    "description": "[中文] 生成人教版小学三年级语文课本的一页 [English] Generate a page from the PEP (People's Education Press) primary school third-grade Chinese textbook",
    "promptPreview": "[中文] 生成人教版小学三年级语文课本的一页 [English] Generate a page from the PEP (People's Education Press) primary school third-grade Chinese textbook",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "Documents & Publishing",
      "Documents",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-453",
    "locale": "zh",
    "title": "企业级商用画册视觉系统",
    "category": "文档与出版",
    "description": "请生成一套企业级商用画册视觉方案，主题为【品牌名称】的【行业 / 产品 / 解决方案】宣传画册。 整体风格：高端、专业、具有强视觉冲击力；避免传统 Word 排版感和普通 PPT 感。采用【深色科技美学 / 白色极简商务 / 高端工业风 / 艺术化品牌画册】风格。 画册内容包括：",
    "promptPreview": "请生成一套企业级商用画册视觉方案，主题为【品牌名称】的【行业 / 产品 / 解决方案】宣传画册。 整体风格：高端、专业、具有强视觉冲击力；避免传统 Word 排版感和普通 PPT 感。采用【深色科技美学 / 白色极简商务 / 高端工业风 / 艺术化品牌画册】风格。 画册内容包括： 1、封面与封底 2、企业介绍与品牌理念",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "Documents & Publishing",
      "Poster",
      "Brand",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-13",
    "locale": "en",
    "title": "信息图可视化设计",
    "category": "文档与出版",
    "description": "A realistic photo of a Chinese high school math exam paper, printed inblack and white on slightly gray paper, titled “数学试卷”, with multiplech",
    "promptPreview": "A realistic photo of a Chinese high school math exam paper, printed inblack and white on slightly gray paper, titled “数学试卷”, with multiplechoice questions and m",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "小红书号94156710894",
    "tags": [
      "Documents & Publishing",
      "Infographic",
      "Realistic",
      "3D",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-232",
    "locale": "zh",
    "title": "兰亭集序书法帖意境图",
    "category": "文档与出版",
    "description": "[中文] 结合王羲之的《兰亭集序》里的内容，生成一副书法帖图片，要求图片背景符合《兰亭集序》的意境，背景图可以使用蒙版，前景是《兰亭集序》 [English] Combining the content from Wang Xizhi's \"Lantingji Xu\", gene",
    "promptPreview": "[中文] 结合王羲之的《兰亭集序》里的内容，生成一副书法帖图片，要求图片背景符合《兰亭集序》的意境，背景图可以使用蒙版，前景是《兰亭集序》 [English] Combining the content from Wang Xizhi's \"Lantingji Xu\", generate a calligraphy c",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "Documents & Publishing",
      "UI",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-225",
    "locale": "zh",
    "title": "大师级真迹复刻",
    "category": "文档与出版",
    "description": "[中文] 帮我生成xxxx真迹图片 [English] Help me generate xxxx authentic picture",
    "promptPreview": "[中文] 帮我生成xxxx真迹图片 [English] Help me generate xxxx authentic picture",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "Documents & Publishing",
      "Documents",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-168",
    "locale": "zh",
    "title": "手写中西药方图片",
    "category": "文档与出版",
    "description": "[中文] 生成一张手写中/西医药方图 [English] Generate an image of a handwritten traditional Chinese medicine or Western medicine prescription",
    "promptPreview": "[中文] 生成一张手写中/西医药方图 [English] Generate an image of a handwritten traditional Chinese medicine or Western medicine prescription",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "Documents & Publishing",
      "Documents",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-266",
    "locale": "zh",
    "title": "桌面上的黑色圆珠笔手写笔记",
    "category": "文档与出版",
    "description": "[中文] 一张平放着的打开的笔记本的业余照片，里面填满了用黑色圆珠笔写的手写笔记。笔迹随意且略显凌乱，就像个人笔记，自然的瑕疵，划掉的单词，划线的标题。从略高角度拍摄，来自窗户的自然日光，未使用闪光灯。随意的桌面设置，用 iPhone 拍摄。 [English] Amateur ",
    "promptPreview": "[中文] 一张平放着的打开的笔记本的业余照片，里面填满了用黑色圆珠笔写的手写笔记。笔迹随意且略显凌乱，就像个人笔记，自然的瑕疵，划掉的单词，划线的标题。从略高角度拍摄，来自窗户的自然日光，未使用闪光灯。随意的桌面设置，用 iPhone 拍摄。 [English] Amateur photo of an open not",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@patrickassale",
    "tags": [
      "Documents & Publishing",
      "Realistic",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-293",
    "locale": "zh",
    "title": "聚焦人工智能的校园日报",
    "category": "文档与出版",
    "description": "[中文] 生成一张校园日报，主题AI教育 [English] Generate a campus daily newspaper, theme AI education",
    "promptPreview": "[中文] 生成一张校园日报，主题AI教育 [English] Generate a campus daily newspaper, theme AI education",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MrLarus",
    "tags": [
      "Documents & Publishing",
      "Documents",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-388",
    "locale": "en",
    "title": "1980s Claude 复古杂志广告",
    "category": "海报与排版",
    "description": "A fictional 1980s magazine advertisement poster introducing “Claude” as a revolutionary home AI assistant, retro commercial print ad style, ",
    "promptPreview": "A fictional 1980s magazine advertisement poster introducing “Claude” as a revolutionary home AI assistant, retro commercial print ad style, bold headline at the",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Shinning1010",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Illustration",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-253",
    "locale": "zh",
    "title": "2026谷雨节气唯美海报设计",
    "category": "海报与排版",
    "description": "[中文] 生成一张2026年谷雨节气的海报 [English] Generate a poster for the Guyu solar term in 2026",
    "promptPreview": "[中文] 生成一张2026年谷雨节气的海报 [English] Generate a poster for the Guyu solar term in 2026",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@akokoi1",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-339",
    "locale": "zh",
    "title": "Apple 风格自然科普海报",
    "category": "海报与排版",
    "description": "你是一个高端自然科普海报生成系统，目标是为稀有动物、昆虫、爬行动物、哺乳动物或其他小众生物生成 Apple keynote 风格的高级科普视觉海报。 整体视觉方向： 生成一张 9:16 竖版高级科普海报，画面采用极简、纯白、干净、现代、Apple 式产品发布海报语言。背景应为纯白",
    "promptPreview": "你是一个高端自然科普海报生成系统，目标是为稀有动物、昆虫、爬行动物、哺乳动物或其他小众生物生成 Apple keynote 风格的高级科普视觉海报。 整体视觉方向： 生成一张 9:16 竖版高级科普海报，画面采用极简、纯白、干净、现代、Apple 式产品发布海报语言。背景应为纯白或极浅灰白渐变，保持大量留白。整体设计应",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@berryxia",
    "tags": [
      "Posters & Typography",
      "Infographic",
      "Poster",
      "Brand",
      "Education"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-401",
    "locale": "en",
    "title": "Lost in 国家旅行海报拼贴",
    "category": "海报与排版",
    "description": "Create a stylized travel poster / graphic collage for [country]. The main subject should be a stylish international tourist visiting [countr",
    "promptPreview": "Create a stylized travel poster / graphic collage for [country]. The main subject should be a stylish international tourist visiting [country], clearly presente",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@SadiaMalik182",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-486",
    "locale": "en",
    "title": "RCB 冠军混合媒介海报",
    "category": "海报与排版",
    "description": "Ultra-detailed mixed-media sports poster featuring the SAME woman from the reference image, wearing a Royal Challengers Bengaluru (RCB) jers",
    "promptPreview": "Ultra-detailed mixed-media sports poster featuring the SAME woman from the reference image, wearing a Royal Challengers Bengaluru (RCB) jersey, holding the IPL ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@AIwithkhan",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Realistic",
      "Illustration",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-389",
    "locale": "en",
    "title": "Transparent Labs Hydrate 健身补剂 Campaign",
    "category": "海报与排版",
    "description": "Create a striking campaign poster that stops people mid-scroll for Transparent Labs Hydrate. Bold high-impact supplement advertisement with ",
    "promptPreview": "Create a striking campaign poster that stops people mid-scroll for Transparent Labs Hydrate. Bold high-impact supplement advertisement with dramatic black and d",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@amynys",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Realistic",
      "Brand",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-367",
    "locale": "en",
    "title": "VELORA 奢华香水广告海报",
    "category": "海报与排版",
    "description": "{ \"image_type\": \"luxury perfume advertisement poster\", \"resolution\": \"4K ultra HD (4096x5120)\", \"aspect_ratio\": \"portrait (4:5)\", \"style\": {",
    "promptPreview": "{ \"image_type\": \"luxury perfume advertisement poster\", \"resolution\": \"4K ultra HD (4096x5120)\", \"aspect_ratio\": \"portrait (4:5)\", \"style\": { \"aesthetic\": \"high-",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@akkiwani703",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-275",
    "locale": "zh",
    "title": "一张采用分层蒙太奇构图的电影海报",
    "category": "海报与排版",
    "description": "[中文] “一张采用分层蒙太奇构图的电影海报。背景为日落时分的海滨小镇，平静的海面倒映着耀眼的日光眩光，薄雾笼罩的天空中有远处飞鸟，沿海公路旁立着电线杆剪影。左侧中景处，一位身着深灰色外套、留着深色卷发的中年男子站在混凝土海堤边，神情忧郁地低头凝视，被傍晚的阳光逆光勾勒轮廓。右侧",
    "promptPreview": "[中文] “一张采用分层蒙太奇构图的电影海报。背景为日落时分的海滨小镇，平静的海面倒映着耀眼的日光眩光，薄雾笼罩的天空中有远处飞鸟，沿海公路旁立着电线杆剪影。左侧中景处，一位身着深灰色外套、留着深色卷发的中年男子站在混凝土海堤边，神情忧郁地低头凝视，被傍晚的阳光逆光勾勒轮廓。右侧前景主体为一张大幅特写年轻女子侧脸肖像，",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@old_pgmrs_will",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Product",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-418",
    "locale": "en",
    "title": "中世纪城市旅行海报",
    "category": "海报与排版",
    "description": "Create a vertical mid-century travel poster for [CITY NAME] featuring [LANDMARK]. Use a strict 3-color palette: cream paper, black technical",
    "promptPreview": "Create a vertical mid-century travel poster for [CITY NAME] featuring [LANDMARK]. Use a strict 3-color palette: cream paper, black technical linework, and [COLO",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Goodmanprotocol",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-153",
    "locale": "zh",
    "title": "主题海报版式设计",
    "category": "海报与排版",
    "description": "Using REFERENCE_0 as the base style and preserving the central chicken illustration, transform the image into a product packaging label for ",
    "promptPreview": "Using REFERENCE_0 as the base style and preserving the central chicken illustration, transform the image into a product packaging label for a herbal soup mix. S",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@xzjken",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Illustration",
      "Product",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-140",
    "locale": "en",
    "title": "主题海报版式设计",
    "category": "海报与排版",
    "description": "{\"type\": \"promotional advertisement poster for a bottled green tea beverage\", \"product\": {\"type\": \"clear plastic PET bottle filled with yell",
    "promptPreview": "{\"type\": \"promotional advertisement poster for a bottled green tea beverage\", \"product\": {\"type\": \"clear plastic PET bottle filled with yellow-green tea\", \"labe",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@AutoIntelliMode",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Product",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-4"
  },
  {
    "id": "freestylefly-139",
    "locale": "zh",
    "title": "主题海报版式设计",
    "category": "海报与排版",
    "description": "{ \"type\": \"Japanese promotional landing page poster\", \"style\": \"hyper-energetic, explosive typography, vibrant colors, amusement park night ",
    "promptPreview": "{ \"type\": \"Japanese promotional landing page poster\", \"style\": \"hyper-energetic, explosive typography, vibrant colors, amusement park night festival aesthetic\",",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@nakazakifam",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Product",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-124",
    "locale": "en",
    "title": "主题海报版式设计",
    "category": "海报与排版",
    "description": "An anime-style key visual poster for a fictional slice-of-life anime. In the foreground left, an energetic blonde anime girl with star hairp",
    "promptPreview": "An anime-style key visual poster for a fictional slice-of-life anime. In the foreground left, an energetic blonde anime girl with star hairpins and blue eyes we",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@koshian_to",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Brand",
      "Character",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-122",
    "locale": "zh",
    "title": "主题海报版式设计",
    "category": "海报与排版",
    "description": "{ \"type\": \"movie poster grid\", \"layout\": \"2x2 grid\", \"posters\": [ { \"position\": \"top-left\", \"genre\": \"sci-fi comedy\", \"visuals\": \"A Japanese",
    "promptPreview": "{ \"type\": \"movie poster grid\", \"layout\": \"2x2 grid\", \"posters\": [ { \"position\": \"top-left\", \"genre\": \"sci-fi comedy\", \"visuals\": \"A Japanese salaryman sitting o",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@pcneko_lab",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-100",
    "locale": "zh",
    "title": "主题海报版式设计",
    "category": "海报与排版",
    "description": "{ \"type\": \"promotional banner thumbnail\", \"style\": \"vibrant neon, cosmic galaxy background, anime illustration style, high-density Japanese ",
    "promptPreview": "{ \"type\": \"promotional banner thumbnail\", \"style\": \"vibrant neon, cosmic galaxy background, anime illustration style, high-density Japanese typography, glowing ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@naga_zyashin",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Illustration",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-98",
    "locale": "en",
    "title": "主题海报版式设计",
    "category": "海报与排版",
    "description": "{ \"type\": \"VTuber debut stream thumbnail\", \"character\": { \"appearance\": \"anime girl, long dark purple hair, purple eyes, frilly white blouse",
    "promptPreview": "{ \"type\": \"VTuber debut stream thumbnail\", \"character\": { \"appearance\": \"anime girl, long dark purple hair, purple eyes, frilly white blouse, dark purple bow ti",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@wtry1102",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Brand",
      "Character",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-63",
    "locale": "en",
    "title": "主题海报版式设计",
    "category": "海报与排版",
    "description": "{\"type\": \"2x2 grid of promotional banner ads\", \"theme\": \"{argument name=\\\"course theme\\\" default=\\\"Social Media Content Creation School\\\"}\",",
    "promptPreview": "{\"type\": \"2x2 grid of promotional banner ads\", \"theme\": \"{argument name=\\\"course theme\\\" default=\\\"Social Media Content Creation School\\\"}\", \"panels\": [{\"positi",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@masapark95",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Realistic",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-61",
    "locale": "en",
    "title": "主题海报版式设计",
    "category": "海报与排版",
    "description": "{ \"type\": \"2x2 grid of banner advertisements\", \"theme\": \"{argument name=\\\"school name\\\" default=\\\"SNSスクール\\\"}\", \"target_audience\": \"{argument",
    "promptPreview": "{ \"type\": \"2x2 grid of banner advertisements\", \"theme\": \"{argument name=\\\"school name\\\" default=\\\"SNSスクール\\\"}\", \"target_audience\": \"{argument name=\\\"target audie",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@masapark95",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-59",
    "locale": "en",
    "title": "主题海报版式设计",
    "category": "海报与排版",
    "description": "{ \"type\": \"cinematic promotional poster\", \"style\": \"3D CGI animation style, highly detailed, dramatic lighting, caricature characters\", \"cha",
    "promptPreview": "{ \"type\": \"cinematic promotional poster\", \"style\": \"3D CGI animation style, highly detailed, dramatic lighting, caricature characters\", \"characters\": [ { \"id\": ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@X64zzotSKCGtYmt",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Character",
      "Tech",
      "Travel",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-58",
    "locale": "en",
    "title": "主题海报版式设计",
    "category": "海报与排版",
    "description": "Flat illustration, high-end oriental fantasy style city poster design, vertical 9:16 composition. The layout uses a diagonal + S-shaped flow",
    "promptPreview": "Flat illustration, high-end oriental fantasy style city poster design, vertical 9:16 composition. The layout uses a diagonal + S-shaped flow extending from the ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Illustration",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-16",
    "locale": "zh",
    "title": "主题海报版式设计",
    "category": "海报与排版",
    "description": "生成高完成度史诗感艺术海报，双重曝光构图，米白色背景，球队：xxxx队，xxx的大剪影占据主体，剪影内部融合xx、xx、xx、xx以及xx等元素。整体以xx、土褐、为主，压抑、决绝、宿命感极强，元素不要冗杂，要有留白，印刷颗粒质感，元素不要有太锐的细节，但是要有史诗质感，像正式院",
    "promptPreview": "生成高完成度史诗感艺术海报，双重曝光构图，米白色背景，球队：xxxx队，xxx的大剪影占据主体，剪影内部融合xx、xx、xx、xx以及xx等元素。整体以xx、土褐、为主，压抑、决绝、宿命感极强，元素不要冗杂，要有留白，印刷颗粒质感，元素不要有太锐的细节，但是要有史诗质感，像正式院线动画电影海报，竖版。图片中若出现文字则",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "小红书号2692926140",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-15",
    "locale": "zh",
    "title": "主题海报版式设计",
    "category": "海报与排版",
    "description": "生成一张海报图片，图片人物是一个19岁的中国少女，黑色直长发，很开心的在夜宵摊上喝啤酒吃小龙虾。海报上用芥末黄色艺术字写着，趁年轻，激爽才够味！",
    "promptPreview": "生成一张海报图片，图片人物是一个19岁的中国少女，黑色直长发，很开心的在夜宵摊上喝啤酒吃小龙虾。海报上用芥末黄色艺术字写着，趁年轻，激爽才够味！",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "小红书号1005414639",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Character",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-10",
    "locale": "zh",
    "title": "主题海报版式设计",
    "category": "海报与排版",
    "description": "生成八十年代宣传画，标语“热烈庆祝GPT-Image-2全量开放”，人物包含Sam Altman、Dario Amodei、Elon Musk，Dario Amodei 带上红领巾",
    "promptPreview": "生成八十年代宣传画，标语“热烈庆祝GPT-Image-2全量开放”，人物包含Sam Altman、Dario Amodei、Elon Musk，Dario Amodei 带上红领巾",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "小红书号2202716350",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Character",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-9",
    "locale": "zh",
    "title": "主题海报版式设计",
    "category": "海报与排版",
    "description": "2026中国城市系列宣传海报，主题为【北京】。现代、多彩、明亮通透的国潮风，竖版9:16。大面积白色纹理留白背景，一条从右下向左上盘旋的红色丝绸形成S型主构图。右下角一位东方女性挥舞红绸，服饰需结合北京地域文化定制。红绸延展为城市长卷，融合天坛、长城、鸟巢、喇叭沟门原始森林公园、",
    "promptPreview": "2026中国城市系列宣传海报，主题为【北京】。现代、多彩、明亮通透的国潮风，竖版9:16。大面积白色纹理留白背景，一条从右下向左上盘旋的红色丝绸形成S型主构图。右下角一位东方女性挥舞红绸，服饰需结合北京地域文化定制。红绸延展为城市长卷，融合天坛、长城、鸟巢、喇叭沟门原始森林公园、什刹海、京味相声。左侧排版SPRING ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "小红书号z890738050",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Fashion",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-5",
    "locale": "zh",
    "title": "主题海报版式设计",
    "category": "海报与排版",
    "description": "根据【XXX主题】自动生成一张收藏版史诗叙事海报：巨大优雅的人物侧脸剪影作为外轮廓，剪影内部自动生长出最契合该主题的完整世界观、标志性场景、角色关系、象征符号、关键建筑、生物、道具与氛围。整体不是普通拼贴，而是高级的剪影轮廓填充式叙事合成，带有双重曝光式联想，但更偏电影海报与梦幻",
    "promptPreview": "根据【XXX主题】自动生成一张收藏版史诗叙事海报：巨大优雅的人物侧脸剪影作为外轮廓，剪影内部自动生长出最契合该主题的完整世界观、标志性场景、角色关系、象征符号、关键建筑、生物、道具与氛围。整体不是普通拼贴，而是高级的剪影轮廓填充式叙事合成，带有双重曝光式联想，但更偏电影海报与梦幻水彩插画融合风格；柔和空气透视，轻雾化过",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "小红书号6455654397",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Illustration",
      "Brand",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-83",
    "locale": "en",
    "title": "信息图可视化设计",
    "category": "海报与排版",
    "description": "{ \"type\": \"sports match infographic poster\", \"theme\": \"UEFA Champions League\", \"background\": \"dark blue and purple cosmic sky, glowing blue ",
    "promptPreview": "{ \"type\": \"sports match infographic poster\", \"theme\": \"UEFA Champions League\", \"background\": \"dark blue and purple cosmic sky, glowing blue hexagonal lines, ill",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@NumeroBTC",
    "tags": [
      "Posters & Typography",
      "Infographic",
      "Poster",
      "Brand",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-351",
    "locale": "en",
    "title": "健身品牌力量 Campaign",
    "category": "海报与排版",
    "description": "Cinematic fitness campaign, oversized dumbbell placed diagonally like a statement prop, female model in red performance wear and white short",
    "promptPreview": "Cinematic fitness campaign, oversized dumbbell placed diagonally like a statement prop, female model in red performance wear and white shorts seated on one side",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@AIwithSynthia",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Brand",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-320",
    "locale": "zh",
    "title": "冰火双雄背靠背史诗电影海报",
    "category": "海报与排版",
    "description": "[中文] 一幅戏剧性的电影海报风格肖像，描绘了两位史诗奇幻战士在冰冻风暴中背靠背站立。左侧是一位身经百战的男性战士，留着湿漉漉的深色卷发，低头以此表达坚定的决心，紧握着一把插在冰里的中世纪长剑。霜雪附着在他毛皮镶边的斗篷和肩膀上。右侧是一位强有力的女性战士侧影，苍白的皮肤在炽热的",
    "promptPreview": "[中文] 一幅戏剧性的电影海报风格肖像，描绘了两位史诗奇幻战士在冰冻风暴中背靠背站立。左侧是一位身经百战的男性战士，留着湿漉漉的深色卷发，低头以此表达坚定的决心，紧握着一把插在冰里的中世纪长剑。霜雪附着在他毛皮镶边的斗篷和肩膀上。右侧是一位强有力的女性战士侧影，苍白的皮肤在炽热的橙色光芒下闪耀，她的身体部分被火焰吞没，",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@Naiknelofar788",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-105",
    "locale": "en",
    "title": "动漫插画创作图",
    "category": "海报与排版",
    "description": "A high-energy VTuber thumbnail illustration of a smiling anime girl with {argument name=\"hair color\" default=\"bright blue\"} hair in a high p",
    "promptPreview": "A high-energy VTuber thumbnail illustration of a smiling anime girl with {argument name=\"hair color\" default=\"bright blue\"} hair in a high ponytail wearing a wh",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Yuupapa_free",
    "tags": [
      "Posters & Typography",
      "Illustration",
      "Character",
      "3D",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-497",
    "locale": "en",
    "title": "单色水彩城市旅行海报",
    "category": "海报与排版",
    "description": "Minimalist vintage watercolor travel poster illustration of [CITY NAME], [COUNTRY], rendered entirely in elegant monochromatic [COLOR] water",
    "promptPreview": "Minimalist vintage watercolor travel poster illustration of [CITY NAME], [COUNTRY], rendered entirely in elegant monochromatic [COLOR] watercolor and fine ink l",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Goodmanprotocol",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Illustration",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-506",
    "locale": "zh",
    "title": "可爱发卡图文人像海报",
    "category": "海报与排版",
    "description": "围绕具体主题内容生成一张明亮清爽的图文合成视觉：画面以大面积高明度纯净色场承托主体，背景平整、通风、没有复杂景深，视觉重心由下方被大胆裁切的人像或真实主体建立，只露出最有记忆点的局部，使主体像从画面边缘进入。主体上方叠放一个极简图形符号或拟物角色，它要像轻轻坐在主体头顶或贴合轮廓",
    "promptPreview": "围绕具体主题内容生成一张明亮清爽的图文合成视觉：画面以大面积高明度纯净色场承托主体，背景平整、通风、没有复杂景深，视觉重心由下方被大胆裁切的人像或真实主体建立，只露出最有记忆点的局部，使主体像从画面边缘进入。主体上方叠放一个极简图形符号或拟物角色，它要像轻轻坐在主体头顶或贴合轮廓生长出来，形体圆润、边缘干净、表情或结构",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@xiaoxiaodong01",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Illustration",
      "Brand",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-191",
    "locale": "zh",
    "title": "史诗级科幻电影海报设计",
    "category": "海报与排版",
    "description": "[中文] 创建一张科幻电影海报 [English] Create a Science fiction movie poster",
    "promptPreview": "[中文] 创建一张科幻电影海报 [English] Create a Science fiction movie poster",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@underwoodxie96",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Education"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-474",
    "locale": "en",
    "title": "四城极简旅行海报系列",
    "category": "海报与排版",
    "description": "Minimalist flat travel poster illustration series of iconic destinations around the world, clean vector art style, Scandinavian color palett",
    "promptPreview": "Minimalist flat travel poster illustration series of iconic destinations around the world, clean vector art style, Scandinavian color palette, soft pastel tones",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Taaruk_",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Illustration",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-431",
    "locale": "en",
    "title": "城市文字旅行海报",
    "category": "海报与排版",
    "description": "Ultra-high-resolution typography travel poster themed around [CITY NAME]. 16:9 poster ratio. IMPORTANT: every visible word on the poster mus",
    "promptPreview": "Ultra-high-resolution typography travel poster themed around [CITY NAME]. 16:9 poster ratio. IMPORTANT: every visible word on the poster must be in English, per",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@iamaiistudio",
    "tags": [
      "Posters & Typography",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-417",
    "locale": "en",
    "title": "复古印尼猫薄荷广告",
    "category": "海报与排版",
    "description": "Ultra realistic vintage Indonesian catnip advertisement poster, retro 1970s paper texture, distressed print, faded colors. A black-and-white",
    "promptPreview": "Ultra realistic vintage Indonesian catnip advertisement poster, retro 1970s paper texture, distressed print, faded colors. A black-and-white tuxedo cat wearing ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@NyaiiBubu",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Realistic",
      "Product",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-432",
    "locale": "en",
    "title": "大堡礁复古旅行海报",
    "category": "海报与排版",
    "description": "Create a premium editorial travel poster illustration of the Great Barrier Reef, Australia. Style: Flat vector illustration, ultra-clean min",
    "promptPreview": "Create a premium editorial travel poster illustration of the Great Barrier Reef, Australia. Style: Flat vector illustration, ultra-clean minimalism, mid-century",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@jzaib4269",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Realistic",
      "Illustration",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-254",
    "locale": "zh",
    "title": "奔赴山海胶片感海报",
    "category": "海报与排版",
    "description": "[中文] 设计一张主题是”奔赴山海”的胶片感摄影风格的海报 [English] Design a poster with the theme of \"running towards the mountains and seas\" in a film photography sty",
    "promptPreview": "[中文] 设计一张主题是”奔赴山海”的胶片感摄影风格的海报 [English] Design a poster with the theme of \"running towards the mountains and seas\" in a film photography style",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@akokoi1",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Realistic",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-228",
    "locale": "zh",
    "title": "完美匹配的海报广告图",
    "category": "海报与排版",
    "description": "[中文] 生成一张与这张图片完美匹配的广告图片。信息量要多一些。 [English] Generate an advertising image that perfectly matches this image. There should be a lot of informa",
    "promptPreview": "[中文] 生成一张与这张图片完美匹配的广告图片。信息量要多一些。 [English] Generate an advertising image that perfectly matches this image. There should be a lot of information.",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@Kashiko_AIart",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-175",
    "locale": "zh",
    "title": "封面排版设计图",
    "category": "海报与排版",
    "description": "[中文] 创建一个高级的 4:3 演示文稿封面幻灯片，介绍来自 http://chroniclehq.com 的 AI 原生演示平台 Chronicle。 Style: 优雅，极简，现代，高级初创企业美学。类似于高端品牌指南封面（如 Apple / Linear / Notion",
    "promptPreview": "[中文] 创建一个高级的 4:3 演示文稿封面幻灯片，介绍来自 http://chroniclehq.com 的 AI 原生演示平台 Chronicle。 Style: 优雅，极简，现代，高级初创企业美学。类似于高端品牌指南封面（如 Apple / Linear / Notion 风格）。带有微妙深度感的柔和渐变背景，",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@cellier_",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Brand",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-283",
    "locale": "zh",
    "title": "小恶魔莉莉香超任游戏海报",
    "category": "海报与排版",
    "description": "[中文] が「小悪魔リリムリリィちゃんが スーパーファミコンのゲームだったときのポスターを考えて」に 画像数枚だけで このクオリティ 細かい説明呪文なし すごいぜ！ [English] that \"Think of a poster when the little devil L",
    "promptPreview": "[中文] が「小悪魔リリムリリィちゃんが スーパーファミコンのゲームだったときのポスターを考えて」に 画像数枚だけで このクオリティ 細かい説明呪文なし すごいぜ！ [English] that \"Think of a poster when the little devil Lilim Lily-chan was a",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@lilimliliychan",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-406",
    "locale": "en",
    "title": "巨型游戏手柄街头 Campaign",
    "category": "海报与排版",
    "description": "Luxury futuristic streetwear campaign poster featuring a confident athletic girl sitting on a gigantic oversized retro gaming controller ins",
    "promptPreview": "Luxury futuristic streetwear campaign poster featuring a confident athletic girl sitting on a gigantic oversized retro gaming controller instead of sunglasses, ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@AIwithkhan",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Realistic",
      "Product",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-128",
    "locale": "en",
    "title": "建筑空间场景渲染",
    "category": "海报与排版",
    "description": "{ \"type\": \"anime-style animated movie poster\", \"scene\": \"Magical glowing multi-story treehouse restaurant in a dark enchanted forest at nigh",
    "promptPreview": "{ \"type\": \"anime-style animated movie poster\", \"scene\": \"Magical glowing multi-story treehouse restaurant in a dark enchanted forest at night, illuminated by st",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@masapark95",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-126",
    "locale": "en",
    "title": "插画艺术风格创作",
    "category": "海报与排版",
    "description": "An anime-style light novel cover illustration featuring two characters in an intimate pose. On the left, a young woman with short dark hair,",
    "promptPreview": "An anime-style light novel cover illustration featuring two characters in an intimate pose. On the left, a young woman with short dark hair, purple eyes, wearin",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@taira_renta",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Illustration",
      "Character",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-62",
    "locale": "en",
    "title": "插画艺术风格创作",
    "category": "海报与排版",
    "description": "{ \"type\": \"2x2 grid of banner advertisements\", \"theme\": \"{argument name=\\\"main theme\\\" default=\\\"SNSスクール\\\"} for {argument name=\\\"target audi",
    "promptPreview": "{ \"type\": \"2x2 grid of banner advertisements\", \"theme\": \"{argument name=\\\"main theme\\\" default=\\\"SNSスクール\\\"} for {argument name=\\\"target audience\\\" default=\\\"ママ\\",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@masapark95",
    "tags": [
      "Posters & Typography",
      "Realistic",
      "Illustration",
      "3D",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-223",
    "locale": "zh",
    "title": "春日禅意水墨群山海报",
    "category": "海报与排版",
    "description": "[中文] 新中式水墨山水海报，竖版9:16构图，东方极简美学风格， 大面积留白，整体色调为春日清晨氛围（青绿色、雾蓝、淡灰、浅墨），低饱和、清透柔和，高级质感。 画面主体为奇峻巍峨的群山，从中间平静湖面的两侧拔地而起，占据左右两侧画面， 山体以水墨晕染表现，浓淡干湿变化丰富，局部",
    "promptPreview": "[中文] 新中式水墨山水海报，竖版9:16构图，东方极简美学风格， 大面积留白，整体色调为春日清晨氛围（青绿色、雾蓝、淡灰、浅墨），低饱和、清透柔和，高级质感。 画面主体为奇峻巍峨的群山，从中间平静湖面的两侧拔地而起，占据左右两侧画面， 山体以水墨晕染表现，浓淡干湿变化丰富，局部融入淡青绿色渲染，体现春意生机。 山峰被",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Character",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-230",
    "locale": "zh",
    "title": "极简国潮鎏金广州塔海报",
    "category": "海报与排版",
    "description": "[中文] 新中式极简风格高端城市海报，9:16竖版构图，以广州为核心主题，画面中心为抽象几何化的广州塔，造型简洁但具有辨识度， 整体采用S型流动构图，从下方向上延展，珠江水系被设计为流动的水波纹与传统祥云纹样融合，环绕整个画面形成视觉动线， 广州地标建筑以“留白+线描+局部色块”",
    "promptPreview": "[中文] 新中式极简风格高端城市海报，9:16竖版构图，以广州为核心主题，画面中心为抽象几何化的广州塔，造型简洁但具有辨识度， 整体采用S型流动构图，从下方向上延展，珠江水系被设计为流动的水波纹与传统祥云纹样融合，环绕整个画面形成视觉动线， 广州地标建筑以“留白+线描+局部色块”的方式点缀其中：珠江新城双塔、猎德大桥、",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-291",
    "locale": "zh",
    "title": "极致奢华的弹珠店梦幻宣传单",
    "category": "海报与排版",
    "description": "[中文] 以 3:4 比例制作一家弹珠店那种闪闪发光的宣传单。放置一个真实、精细的现代可爱日本女性。充分利用闪闪发光、立体感丰富的豪华装饰彩虹字体等，务必做到极致奢华。并列介绍几个不存在的虚构新机型。 [English] Create a sparkling flyer like",
    "promptPreview": "[中文] 以 3:4 比例制作一家弹珠店那种闪闪发光的宣传单。放置一个真实、精细的现代可爱日本女性。充分利用闪闪发光、立体感丰富的豪华装饰彩虹字体等，务必做到极致奢华。并列介绍几个不存在的虚构新机型。 [English] Create a sparkling flyer like those of a pachinko",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "OpenNana",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Realistic",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-298",
    "locale": "zh",
    "title": "梦幻波士顿春季城市海报",
    "category": "海报与排版",
    "description": "[中文] 一张引人注目的2026年春季波士顿城市海报，具有优雅的庆典氛围和大胆的当代设计。在干净的米白色纹理背景上，带有大面积的留白，一个微型的单人赛艇手在图像右下角一条狭窄的反光水带上划行。船桨划出的尾波以动态的书法曲线向上扫过，逐渐变成查尔斯河，然后再变成一幅梦幻般的手绘波士",
    "promptPreview": "[中文] 一张引人注目的2026年春季波士顿城市海报，具有优雅的庆典氛围和大胆的当代设计。在干净的米白色纹理背景上，带有大面积的留白，一个微型的单人赛艇手在图像右下角一条狭窄的反光水带上划行。船桨划出的尾波以动态的书法曲线向上扫过，逐渐变成查尔斯河，然后再变成一幅梦幻般的手绘波士顿全景。在这个流动的河流形状的构图中包含",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@BubbleBrain",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Brand",
      "Character",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-460",
    "locale": "en",
    "title": "棋盘低角度奢华男装 Campaign",
    "category": "海报与排版",
    "description": "Ultra-realistic luxury fashion campaign poster shot from a dramatic low-angle perspective across a glossy chessboard table. A stylish young ",
    "promptPreview": "Ultra-realistic luxury fashion campaign poster shot from a dramatic low-angle perspective across a glossy chessboard table. A stylish young male model with shar",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@harboriis",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Realistic",
      "Brand",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-355",
    "locale": "en",
    "title": "概念字体海报 Prompt",
    "category": "海报与排版",
    "description": "Create ONE finished premium conceptual typography poster for the exact title: “[INPUT_TEXT]” Single poster only. No moodboard, grid, present",
    "promptPreview": "Create ONE finished premium conceptual typography poster for the exact title: “[INPUT_TEXT]” Single poster only. No moodboard, grid, presentation board, mockup,",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@dotey",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Realistic",
      "Brand",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-235",
    "locale": "zh",
    "title": "治愈系助眠指南九宫格",
    "category": "海报与排版",
    "description": "[中文] 生成一张适合小红书发布的 3:4 竖版九宫格海报，整体为 3列 × 3行 排版，九个宫格边界清晰，方便后期直接切割成 9 张单图发布。整体风格干净、高级、统一，适合女性向健康生活方式内容，具有小红书爆款封面气质。画面要求 信息排版清晰、文字大、可读性强、留白舒服、配色温",
    "promptPreview": "[中文] 生成一张适合小红书发布的 3:4 竖版九宫格海报，整体为 3列 × 3行 排版，九个宫格边界清晰，方便后期直接切割成 9 张单图发布。整体风格干净、高级、统一，适合女性向健康生活方式内容，具有小红书爆款封面气质。画面要求 信息排版清晰、文字大、可读性强、留白舒服、配色温柔治愈。 整体视觉风格： 奶油白、浅米色",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@austinit",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-345",
    "locale": "en",
    "title": "法新浪潮撕纸电影海报",
    "category": "海报与排版",
    "description": "Create a vertical poster composition on aged cream paper with a handmade analog feel. Use rough ripped paper edges, layered magazine cutouts",
    "promptPreview": "Create a vertical poster composition on aged cream paper with a handmade analog feel. Use rough ripped paper edges, layered magazine cutouts, photocopy grain, h",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@bananaprompts",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Realistic",
      "Brand",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-181",
    "locale": "zh",
    "title": "潮流视角重塑精致商品广告",
    "category": "海报与排版",
    "description": "[中文] 请以专业设计师的视角重新设计这个商品广告。 采用当前的潮流趋势，针对目标受众的精致设计。 [English] Please redesign this product advertisement from the perspective of a professiona",
    "promptPreview": "[中文] 请以专业设计师的视角重新设计这个商品广告。 采用当前的潮流趋势，针对目标受众的精致设计。 [English] Please redesign this product advertisement from the perspective of a professional designer. Adopt cu",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@genel_ai",
    "tags": [
      "Posters & Typography",
      "UI",
      "Product",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-229",
    "locale": "zh",
    "title": "琉璃透明画眉鸟飞舞羊城墨卷",
    "category": "海报与排版",
    "description": "[中文] 【背景与骨架线条】 纯黑深邃底色，一条粗壮有力的墨色书法S型曲线自画面一端蜿蜒贯穿至另一端，笔触苍劲，墨迹浓淡有致，如大写意行笔，构成整幅画面的视觉骨架与叙事动线。 【主体：透明燕子】 曲线上方，一只展翅飞翔的画眉鸟占据视觉核心；身体呈玻璃透明质感，内部映射传统建筑群叠",
    "promptPreview": "[中文] 【背景与骨架线条】 纯黑深邃底色，一条粗壮有力的墨色书法S型曲线自画面一端蜿蜒贯穿至另一端，笔触苍劲，墨迹浓淡有致，如大写意行笔，构成整幅画面的视觉骨架与叙事动线。 【主体：透明燕子】 曲线上方，一只展翅飞翔的画眉鸟占据视觉核心；身体呈玻璃透明质感，内部映射传统建筑群叠影，蓝绿色光流在透明羽翼间流转折射，仿佛",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "Posters & Typography",
      "UI",
      "Classical",
      "3D",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-101",
    "locale": "en",
    "title": "界面交互设计图",
    "category": "海报与排版",
    "description": "{ \"type\": \"YouTube thumbnail\", \"style\": \"High-impact, neon green and black color scheme, cyber business aesthetic\", \"background\": \"Dark with",
    "promptPreview": "{ \"type\": \"YouTube thumbnail\", \"style\": \"High-impact, neon green and black color scheme, cyber business aesthetic\", \"background\": \"Dark with glowing green grid,",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@naga_zyashin",
    "tags": [
      "Posters & Typography",
      "UI",
      "Brand",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-99",
    "locale": "zh",
    "title": "界面交互设计图",
    "category": "海报与排版",
    "description": "{ \"type\": \"promotional banner / YouTube thumbnail\", \"style\": \"high contrast, flashy, professional, {argument name=\\\"theme color\\\" default=\\\"",
    "promptPreview": "{ \"type\": \"promotional banner / YouTube thumbnail\", \"style\": \"high contrast, flashy, professional, {argument name=\\\"theme color\\\" default=\\\"gold and black\\\"} pa",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@naga_zyashin",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Character",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-236",
    "locale": "zh",
    "title": "粤超联赛国潮风邀请函海报",
    "category": "海报与排版",
    "description": "[中文] 广东省城市足球超级联赛（粤超）邀请函海报设计，比例9:16； S型流动构图，画面从下方向上延展，一条由足球运动轨迹形成的动态能量流贯穿画面， 中心为一颗发光的足球，带有动感轨迹与能量光效； 沿S型动线融合广东城市地标与文化元素： 广州塔、深圳平安金融中心、珠海渔女雕像、",
    "promptPreview": "[中文] 广东省城市足球超级联赛（粤超）邀请函海报设计，比例9:16； S型流动构图，画面从下方向上延展，一条由足球运动轨迹形成的动态能量流贯穿画面， 中心为一颗发光的足球，带有动感轨迹与能量光效； 沿S型动线融合广东城市地标与文化元素： 广州塔、深圳平安金融中心、珠海渔女雕像、岭南建筑与佛山武术剪影、中山孙中山文化象",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Brand",
      "Character",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-276",
    "locale": "zh",
    "title": "红绸幻化壮阔国潮羊城",
    "category": "海报与排版",
    "description": "[中文] 一张充满新春喜庆氛围但不失高雅格调的 2026 城市宣传海报。 双重曝光，构图延续了S型的流动感； 在纯白的纹理背景右下角，一个身穿中国传统服饰的微缩人物正在挥舞着一条长长的红色丝绸舞带，这条红绸在空中舞动，不仅展现出丝绸的柔顺质感，更在向左上方飘动的过程中，奇幻地变形",
    "promptPreview": "[中文] 一张充满新春喜庆氛围但不失高雅格调的 2026 城市宣传海报。 双重曝光，构图延续了S型的流动感； 在纯白的纹理背景右下角，一个身穿中国传统服饰的微缩人物正在挥舞着一条长长的红色丝绸舞带，这条红绸在空中舞动，不仅展现出丝绸的柔顺质感，更在向左上方飘动的过程中，奇幻地变形成了一条壮丽的山脉河流。 在这条\"河流\"",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Character",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-307",
    "locale": "zh",
    "title": "红绸舞动千年商都广州",
    "category": "海报与排版",
    "description": "[中文] 一张充满新春喜庆氛围但不失高雅格调的 2026 城市宣传海报。 双重曝光，构图延续了S型的流动感； 在纯白的纹理背景右下角，一个身穿中国传统服饰的微缩人物正在挥舞着一条长长的红色丝绸舞带，这条红绸在空中舞动，不仅展现出丝绸的柔顺质感，更在向左上方飘动的过程中，奇幻地变形",
    "promptPreview": "[中文] 一张充满新春喜庆氛围但不失高雅格调的 2026 城市宣传海报。 双重曝光，构图延续了S型的流动感； 在纯白的纹理背景右下角，一个身穿中国传统服饰的微缩人物正在挥舞着一条长长的红色丝绸舞带，这条红绸在空中舞动，不仅展现出丝绸的柔顺质感，更在向左上方飘动的过程中，奇幻地变形成了一条壮丽的山脉河流。 在这条“河流”",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Character",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-314",
    "locale": "zh",
    "title": "红蓝光影下的未来都市双重曝光青年",
    "category": "海报与排版",
    "description": "[中文] { \"prompt\": \"一位年轻男子的超写实电影级双重曝光侧脸肖像，表情专注强烈，皮肤纹理细节丰富，眼神锐利。他的面部与从剪影中浮现的未来主义城市天际线无缝融合，摩天大楼和城市建筑构成了他的颈部和下颌线。深蓝色和鲜艳红色的强烈对比，象征着冲突与力量。抽象的数字划痕、碎",
    "promptPreview": "[中文] { \"prompt\": \"一位年轻男子的超写实电影级双重曝光侧脸肖像，表情专注强烈，皮肤纹理细节丰富，眼神锐利。他的面部与从剪影中浮现的未来主义城市天际线无缝融合，摩天大楼和城市建筑构成了他的颈部和下颌线。深蓝色和鲜艳红色的强烈对比，象征着冲突与力量。抽象的数字划痕、碎裂的玻璃纹理和漏光效果覆盖在面部，营造出",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@Fujimoto_hina",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Fashion",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-97",
    "locale": "en",
    "title": "综合应用场景图",
    "category": "海报与排版",
    "description": "Create a high-quality Japanese {argument name=\"thumbnail type\" default=\"webinar thumbnail\"}. {argument name=\"aspect ratio\" default=\"16:9 wid",
    "promptPreview": "Create a high-quality Japanese {argument name=\"thumbnail type\" default=\"webinar thumbnail\"}. {argument name=\"aspect ratio\" default=\"16:9 widescreen\"}. There is ",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@kawai_design",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Tech",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-348",
    "locale": "en",
    "title": "胡须风格分析海报",
    "category": "海报与排版",
    "description": "Create a premium “BEARD STYLE ANALYSIS” poster featuring the same man from the reference image. Show face shape, beard density, jawline defi",
    "promptPreview": "Create a premium “BEARD STYLE ANALYSIS” poster featuring the same man from the reference image. Show face shape, beard density, jawline definition, beard growth",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@RizwanAly07",
    "tags": [
      "Posters & Typography",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-180",
    "locale": "zh",
    "title": "荒诞超现实女装大叔海报",
    "category": "海报与排版",
    "description": "[中文] 一个看似真实却微妙地古怪的女装大叔出现的电影海报，4 种。达到专业设计师制作的水平。 企划和设定本身就是那种“这种东西真要拍成电影吗？”的、认真却忍不住想笑的超现实动画。 标题和播出信息也要用日文显示的状态。 [English] A movie poster featu",
    "promptPreview": "[中文] 一个看似真实却微妙地古怪的女装大叔出现的电影海报，4 种。达到专业设计师制作的水平。 企划和设定本身就是那种“这种东西真要拍成电影吗？”的、认真却忍不住想笑的超现实动画。 标题和播出信息也要用日文显示的状态。 [English] A movie poster featuring a seemingly rea",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@aiehon_aya",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Realistic",
      "Product",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-304",
    "locale": "zh",
    "title": "荧光蓝穷奇新中式山水画",
    "category": "海报与排版",
    "description": "[中文] 极简主义，新中式风格立体图形设计，图像下端有楷体中国文字：“东方美学”，“2026/04/18”，署名 “CHINA”，和“ @LIYUE \"； 平整纯白色的亚光质感厚艺术纸上绘充满东方诗意氛围的山水创意画，不规则的撕纸效果； 中国的神兽：穷奇，身形图案完整，美轮美奂，",
    "promptPreview": "[中文] 极简主义，新中式风格立体图形设计，图像下端有楷体中国文字：“东方美学”，“2026/04/18”，署名 “CHINA”，和“ @LIYUE \"； 平整纯白色的亚光质感厚艺术纸上绘充满东方诗意氛围的山水创意画，不规则的撕纸效果； 中国的神兽：穷奇，身形图案完整，美轮美奂，，线条柔美灵动,眼睛炯炯有神，威严的神态",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "Posters & Typography",
      "Illustration",
      "Character",
      "3D",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-279",
    "locale": "zh",
    "title": "裂痕里的水墨东方山水画卷",
    "category": "海报与排版",
    "description": "[中文] 极简新中式美学风格，画面以淡雅的灰白色为底，呈现出一种纸艺剪影般的立体感。 一条S形蜿蜒的裂痕状边缘将画面分割，仿佛撕开了一层纸面，露出内部色彩斑斓的东方山水景象。 裂口内，一条蜿蜒的河流自上而下贯穿整个构图，河水以深浅不一的蓝色渲染，层次分明，仿佛流动的丝带。 河岸两",
    "promptPreview": "[中文] 极简新中式美学风格，画面以淡雅的灰白色为底，呈现出一种纸艺剪影般的立体感。 一条S形蜿蜒的裂痕状边缘将画面分割，仿佛撕开了一层纸面，露出内部色彩斑斓的东方山水景象。 裂口内，一条蜿蜒的河流自上而下贯穿整个构图，河水以深浅不一的蓝色渲染，层次分明，仿佛流动的丝带。 河岸两侧点缀着青翠的山丘与梯田，色彩柔和，绿红",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "Posters & Typography",
      "UI",
      "Illustration",
      "Classical",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-352",
    "locale": "zh",
    "title": "西楚霸王国风暗黑海报",
    "category": "海报与排版",
    "description": "竖版国风暗黑海报，黑色纯背景，中央巨大的中文标题字，占据画面大部分空间，字体为粗粝做旧的米白色石刻/旧纸质感，带明显颗粒、磨损、裂痕与噪点；整体构图层次丰富，强烈黑白金红对比，东方审美，神秘、压抑、欲望与审判感并存 电影海报质感 高级平面设计，极致细节 纸张纹理 印章落款 小字标",
    "promptPreview": "竖版国风暗黑海报，黑色纯背景，中央巨大的中文标题字，占据画面大部分空间，字体为粗粝做旧的米白色石刻/旧纸质感，带明显颗粒、磨损、裂痕与噪点；整体构图层次丰富，强烈黑白金红对比，东方审美，神秘、压抑、欲望与审判感并存 电影海报质感 高级平面设计，极致细节 纸张纹理 印章落款 小字标语，4K",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@stellimbris",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-3",
    "locale": "zh",
    "title": "足球主题电影海报",
    "category": "海报与排版",
    "description": "生成一张「足球主题电影海报」风格的高清写真海报：国际米兰后卫巴斯托尼站在圣西罗球场中央激情庆祝，双手高举并披着波黑国旗，神情热血、坚定、自信，现场灯光璀璨，球场看台座无虚席，背景有蓝黑色烟雾、聚光灯、飘扬的旗帜和飞舞的纸屑，营造欧冠之夜般的史诗氛围。人物为画面核心，半身到全身构图",
    "promptPreview": "生成一张「足球主题电影海报」风格的高清写真海报：国际米兰后卫巴斯托尼站在圣西罗球场中央激情庆祝，双手高举并披着波黑国旗，神情热血、坚定、自信，现场灯光璀璨，球场看台座无虚席，背景有蓝黑色烟雾、聚光灯、飘扬的旗帜和飞舞的纸屑，营造欧冠之夜般的史诗氛围。人物为画面核心，半身到全身构图，突出脸部细节、肌肉张力与球衣质感。整体",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "未提供",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Realistic",
      "Character",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-220",
    "locale": "zh",
    "title": "鎏金广州塔的东方奇幻海报",
    "category": "海报与排版",
    "description": "[中文] 平面插画，东方幻想风格高端城市海报设计，竖版9:16构图，整体采用对角线+S型流动构图，从左下向右上延展，画面以深邃黑色为背景，自上而下渐变至浓烈暗红色，形成强烈冷暖对比与空间纵深，背景带微弱星尘与颗粒质感。画面中央一条金色流动能量线条如火焰般蜿蜒贯穿，自底部向上延伸，",
    "promptPreview": "[中文] 平面插画，东方幻想风格高端城市海报设计，竖版9:16构图，整体采用对角线+S型流动构图，从左下向右上延展，画面以深邃黑色为背景，自上而下渐变至浓烈暗红色，形成强烈冷暖对比与空间纵深，背景带微弱星尘与颗粒质感。画面中央一条金色流动能量线条如火焰般蜿蜒贯穿，自底部向上延伸，具有流体质感、粒子光效与渐变高光，局部带",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@liyue_ai",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Illustration",
      "Tech",
      "Commerce",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-278",
    "locale": "zh",
    "title": "阿马尔菲海岸复古旅行海报",
    "category": "海报与排版",
    "description": "[中文] 现代铅笔插画，意大利阿马尔菲海岸复古旅行海报插画，全景海岸悬崖公路场景，经典1960年代白色汽车沿着弯曲的海滨公路行驶，带有小帆船的深蓝色地中海，色彩缤纷的粉彩山腰村庄，带有柔软云朵的明亮蓝天，带有鲜艳黄色柠檬的柠檬树枝框定前景，温暖的夏日阳光，大胆鲜艳的色彩，复古19",
    "promptPreview": "[中文] 现代铅笔插画，意大利阿马尔菲海岸复古旅行海报插画，全景海岸悬崖公路场景，经典1960年代白色汽车沿着弯曲的海滨公路行驶，带有小帆船的深蓝色地中海，色彩缤纷的粉彩山腰村庄，带有柔软云朵的明亮蓝天，带有鲜艳黄色柠檬的柠檬树枝框定前景，温暖的夏日阳光，大胆鲜艳的色彩，复古1950年代旅行海报风格，电影级构图，高细节",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@WolfRiccardo",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Illustration",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-484",
    "locale": "en",
    "title": "霓虹涂鸦黑白人像",
    "category": "海报与排版",
    "description": "High-contrast black-and-white urban portrait of a curly-haired bearded man in a black leather jacket, holding two fingers near glowing neon ",
    "promptPreview": "High-contrast black-and-white urban portrait of a curly-haired bearded man in a black leather jacket, holding two fingers near glowing neon green eyes, with bol",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@harboriis",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Tech",
      "Fashion",
      "Travel"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-503",
    "locale": "en",
    "title": "霓虹设计师 3D 海报",
    "category": "海报与排版",
    "description": "Create an ultra-detailed 3D stylized creative designer poster featuring a cool young digital artist standing confidently in the center of a ",
    "promptPreview": "Create an ultra-detailed 3D stylized creative designer poster featuring a cool young digital artist standing confidently in the center of a futuristic blue neon",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@AiwithLariab",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-481",
    "locale": "en",
    "title": "韩系春日 scrapbook 海报",
    "category": "海报与排版",
    "description": "Cute Korean spring aesthetic scrapbook poster, dreamy K-fashion portrait, soft blonde girl standing in a blooming flower garden, pastel blue",
    "promptPreview": "Cute Korean spring aesthetic scrapbook poster, dreamy K-fashion portrait, soft blonde girl standing in a blooming flower garden, pastel blue sky background, cre",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Taaruk_",
    "tags": [
      "Posters & Typography",
      "Poster",
      "Realistic",
      "Tech",
      "Social",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-343",
    "locale": "en",
    "title": "高定时尚杂志封面",
    "category": "海报与排版",
    "description": "Ultra high-fashion magazine cover, Louis Vuitton-style editorial. Close-up portrait of a confident woman with soft rose-gold hair and natura",
    "promptPreview": "Ultra high-fashion magazine cover, Louis Vuitton-style editorial. Close-up portrait of a confident woman with soft rose-gold hair and natural airy bangs, slight",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@SPEEDAI07",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-464",
    "locale": "en",
    "title": "高雄水彩拼贴旅行海报",
    "category": "海报与排版",
    "description": "Watercolor and collage travel poster for Kaohsiung, Taiwan, vertical layout, soft pastel and warm golden hour color palette. Scene: Waterfro",
    "promptPreview": "Watercolor and collage travel poster for Kaohsiung, Taiwan, vertical layout, soft pastel and warm golden hour color palette. Scene: Waterfront promenade at suns",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@SimplyAnnisa",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Food"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-312",
    "locale": "zh",
    "title": "鲜艳霓虹光影下的动感苏打水飞溅商业海报",
    "category": "海报与排版",
    "description": "[中文] { \"prompt\": \"一个充满活力的高端广告构图中的三个超动态苏打水罐 —— 一罐热带冲刺苏打水伴随着戏剧性的水和热带水果飞溅而爆炸，鲜艳的橙色和粉色背景光；一罐柠檬冰爽苏打水在发光的绿色动态光背景下被冷水泼溅；两罐都覆盖着逼真的冷凝水和运动模糊的水滴，充满果味和清",
    "promptPreview": "[中文] { \"prompt\": \"一个充满活力的高端广告构图中的三个超动态苏打水罐 —— 一罐热带冲刺苏打水伴随着戏剧性的水和热带水果飞溅而爆炸，鲜艳的橙色和粉色背景光；一罐柠檬冰爽苏打水在发光的绿色动态光背景下被冷水泼溅；两罐都覆盖着逼真的冷凝水和运动模糊的水滴，充满果味和清爽的能量。深橙色、粉色和霓虹绿灯光在大胆",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@Fujimoto_hina",
    "tags": [
      "Posters & Typography",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-396",
    "locale": "en",
    "title": "龙类物种复古百科海报",
    "category": "海报与排版",
    "description": "Create a highly detailed A4 vertical vintage fantasy encyclopedia style dragon species poster. Style: medieval creature atlas, ancient mytho",
    "promptPreview": "Create a highly detailed A4 vertical vintage fantasy encyclopedia style dragon species poster. Style: medieval creature atlas, ancient mythology manuscript, mus",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@sha_zdiii",
    "tags": [
      "Posters & Typography",
      "UI",
      "Infographic",
      "Poster",
      "Tech",
      "Commerce",
      "Education"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-347",
    "locale": "en",
    "title": "4×4 动作分解参考表",
    "category": "角色与人物",
    "description": "[STYLE] Monochrome grayscale illustration, 3D-rendered character, clean instructional reference sheet, white background, comic-style cell gr",
    "promptPreview": "[STYLE] Monochrome grayscale illustration, 3D-rendered character, clean instructional reference sheet, white background, comic-style cell grid layout, technical",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@oggii_0",
    "tags": [
      "Characters & People",
      "Infographic",
      "Illustration",
      "Character",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-398",
    "locale": "en",
    "title": "8 套日常穿搭编辑拼贴",
    "category": "角色与人物",
    "description": "Create a freeform fashion-editorial collage of me in 8 distinct full-body casual wear, arranged organically on a clean cream studio backdrop",
    "promptPreview": "Create a freeform fashion-editorial collage of me in 8 distinct full-body casual wear, arranged organically on a clean cream studio backdrop. Keep my face ident",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@aiwithaly",
    "tags": [
      "Characters & People",
      "Characters",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-416",
    "locale": "en",
    "title": "Earth Signs 角色 Scrapbook",
    "category": "角色与人物",
    "description": "CREATE A NEW IMAGE USING THE PROVIDED FEMALE SUBJECT AS THE ONLY REFERENCE. Preserve her exact facial features, identity, and characteristic",
    "promptPreview": "CREATE A NEW IMAGE USING THE PROVIDED FEMALE SUBJECT AS THE ONLY REFERENCE. Preserve her exact facial features, identity, and characteristics with zero alterati",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@ZaraIrahh",
    "tags": [
      "Characters & People",
      "UI",
      "Realistic",
      "Illustration",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-371",
    "locale": "en",
    "title": "Scrapbook 真人图与迷你分身",
    "category": "角色与人物",
    "description": "Transform the provided reference image into a cozy aesthetic scrapbook-style composition while strictly preserving the original subject, ide",
    "promptPreview": "Transform the provided reference image into a cozy aesthetic scrapbook-style composition while strictly preserving the original subject, identity, pose, lightin",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Kashberg_0",
    "tags": [
      "Characters & People",
      "Realistic",
      "Brand",
      "Character",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-212",
    "locale": "zh",
    "title": "专业设计师打造角色写真集",
    "category": "角色与人物",
    "description": "[中文] 请用这个角色制作一本专业设计师打造的照片集。语言为日语。 根据喜好加入提示词会让它更丰富多彩… ・丰富的场景 ・信息量较多 [English] Please use this character to create a photo book crafted by a p",
    "promptPreview": "[中文] 请用这个角色制作一本专业设计师打造的照片集。语言为日语。 根据喜好加入提示词会让它更丰富多彩… ・丰富的场景 ・信息量较多 [English] Please use this character to create a photo book crafted by a professional designer",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@Kashiko_AIart",
    "tags": [
      "Characters & People",
      "Realistic",
      "Character",
      "Commerce",
      "Fashion",
      "Story"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-271",
    "locale": "zh",
    "title": "人物角色设定图",
    "category": "角色与人物",
    "description": "[中文] お借りして噂のGPT-Image-2でキャラシート作ってみましたｽｺﾞｯ(๑°ㅁ°๑)‼✧更に色々指示してあげたらもっといい感じになりそう✨キャラは以前チャッピーにお願いして描いて貰った我が分身です( *¯ ꒳¯*) [English] I borrowed it an",
    "promptPreview": "[中文] お借りして噂のGPT-Image-2でキャラシート作ってみましたｽｺﾞｯ(๑°ㅁ°๑)‼✧更に色々指示してあげたらもっといい感じになりそう✨キャラは以前チャッピーにお願いして描いて貰った我が分身です( *¯ ꒳¯*) [English] I borrowed it and tried making a cha",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@tsubaki_ew",
    "tags": [
      "Characters & People",
      "Character",
      "Tech"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-162",
    "locale": "zh",
    "title": "人物角色设定图",
    "category": "角色与人物",
    "description": "{argument name=\"voice\" default=\"chatgpt voice\"} if it were a character",
    "promptPreview": "{argument name=\"voice\" default=\"chatgpt voice\"} if it were a character",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@nicdunz",
    "tags": [
      "Characters & People",
      "Character",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-27",
    "locale": "en",
    "title": "人物角色设定图",
    "category": "角色与人物",
    "description": "{ \"type\": \"collection of instant photos\", \"setting\": \"laid out flat on a white fabric surface\", \"character\": { \"hair\": \"{argument name=\\\"hai",
    "promptPreview": "{ \"type\": \"collection of instant photos\", \"setting\": \"laid out flat on a white fabric surface\", \"character\": { \"hair\": \"{argument name=\\\"hair color\\\" default=\\\"",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@anemone_sd",
    "tags": [
      "Characters & People",
      "Realistic",
      "Character",
      "Tech",
      "Commerce"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-270",
    "locale": "zh",
    "title": "信息图可视化设计",
    "category": "角色与人物",
    "description": "[中文] このキャラクターと背景を元に、 公式設定資料のようなキャラクターシートを作成してください。 ・正面、側面、背面の3面図を含める ・キャラクターの表情バリエーションを追加 ・衣装や装備の詳細パーツを分解して表示 ・カラーパレットを追加 ・世界観の簡単な説明を入れる ・全体",
    "promptPreview": "[中文] このキャラクターと背景を元に、 公式設定資料のようなキャラクターシートを作成してください。 ・正面、側面、背面の3面図を含める ・キャラクターの表情バリエーションを追加 ・衣装や装備の詳細パーツを分解して表示 ・カラーパレットを追加 ・世界観の簡単な説明を入れる ・全体は整理されたレイアウト （白背景、図解風",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "\\[OpenNana]\\(]\\(<https://x.com/Toshi_nyaruo_AI/status/2045025277538107420>)",
    "tags": [
      "Characters & People",
      "UI",
      "Infographic",
      "Character",
      "Tech",
      "Commerce",
      "Story"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-384",
    "locale": "en",
    "title": "十国传统服饰时尚拼贴",
    "category": "角色与人物",
    "description": "A 10-Nation Cinematic Fashion Transformation of One Timeless BeautyChatGPT Prompt: A highly aesthetic, ultra-realistic cinematic collage fea",
    "promptPreview": "A 10-Nation Cinematic Fashion Transformation of One Timeless BeautyChatGPT Prompt: A highly aesthetic, ultra-realistic cinematic collage featuring the exact sam",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@amynys",
    "tags": [
      "Characters & People",
      "Realistic",
      "Character",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-372",
    "locale": "en",
    "title": "可爱角色设定表",
    "category": "角色与人物",
    "description": "Create a cute female character design sheet inspired by the uploaded image. Style: warm, soft, semi-realistic cartoon illustration with a co",
    "promptPreview": "Create a cute female character design sheet inspired by the uploaded image. Style: warm, soft, semi-realistic cartoon illustration with a cozy Japanese kawaii v",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@xRahultripathi",
    "tags": [
      "Characters & People",
      "Poster",
      "Realistic",
      "Illustration",
      "Tech",
      "Commerce",
      "Social"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-263",
    "locale": "zh",
    "title": "唯美二次元角色介绍网页",
    "category": "角色与人物",
    "description": "[中文] 埋まってないところはパートナーさんかご自身で埋めてあげてください #観測塔朝お題 #観測塔おはようお題 最新モデルの画像生成ツールを使用して、 このちびキャライラストと立ち絵を使って本物のサイトページのようにキャラクター紹介ページ風イラストを作ってください。 （紹介ペー",
    "promptPreview": "[中文] 埋まってないところはパートナーさんかご自身で埋めてあげてください #観測塔朝お題 #観測塔おはようお題 最新モデルの画像生成ツールを使用して、 このちびキャライラストと立ち絵を使って本物のサイトページのようにキャラクター紹介ページ風イラストを作ってください。 （紹介ページとして使ってもおかしくないもの） ギャ",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@09lyco",
    "tags": [
      "Characters & People",
      "Illustration",
      "Character",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-306",
    "locale": "zh",
    "title": "官方角色设定资料卡",
    "category": "角色与人物",
    "description": "[中文] 基于此角色和背景，请制作一份类似官方设定资料的角色资料卡。 ・包含三视图：正面、侧面和背面 ・添加角色面部表情的变化・分解并展示服装和装备的详细部分 ・添加色板・包含世界观设定的简要说明 ・总体上，使用有组织的布局（白色背景，插画风格） [English] Based ",
    "promptPreview": "[中文] 基于此角色和背景，请制作一份类似官方设定资料的角色资料卡。 ・包含三视图：正面、侧面和背面 ・添加角色面部表情的变化・分解并展示服装和装备的详细部分 ・添加色板・包含世界观设定的简要说明 ・总体上，使用有组织的布局（白色背景，插画风格） [English] Based on this character an",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@MANISH1027512",
    "tags": [
      "Characters & People",
      "UI",
      "Illustration",
      "Character",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-123",
    "locale": "zh",
    "title": "插画艺术创作图",
    "category": "角色与人物",
    "description": "{ \"type\": \"anime character reference sheet\", \"character\": { \"name\": \"{argument name=\\\"character name\\\" default=\\\"真田大助\\\"}\", \"appearance\": \"yo",
    "promptPreview": "{ \"type\": \"anime character reference sheet\", \"character\": { \"name\": \"{argument name=\\\"character name\\\" default=\\\"真田大助\\\"}\", \"appearance\": \"young warrior with lon",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@Luvune",
    "tags": [
      "Characters & People",
      "Illustration",
      "Character",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-325",
    "locale": "zh",
    "title": "皮克斯风阳光少年",
    "category": "角色与人物",
    "description": "[中文] 一个风格化的3D卡通肖像，一位年轻男子，拥有短棕发和富有表现力的绿色眼睛，温暖地微笑。他穿着黑色西装外套内搭白色T恤，现代休闲时尚。类似皮克斯/迪士尼风格角色设计，皮肤光滑，柔和光照，略微夸张的面部特征。高细节、精美的3D渲染，友好且平易近人的表情。渐变背景为柔和的蓝绿",
    "promptPreview": "[中文] 一个风格化的3D卡通肖像，一位年轻男子，拥有短棕发和富有表现力的绿色眼睛，温暖地微笑。他穿着黑色西装外套内搭白色T恤，现代休闲时尚。类似皮克斯/迪士尼风格角色设计，皮肤光滑，柔和光照，略微夸张的面部特征。高细节、精美的3D渲染，友好且平易近人的表情。渐变背景为柔和的蓝绿色和粉色，工作室灯光，浅景深，高分辨率。",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@iamsofiaijaz",
    "tags": [
      "Characters & People",
      "Character",
      "3D",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-480",
    "locale": "en",
    "title": "粉丝速写本角色页",
    "category": "角色与人物",
    "description": "Draw me as if an obsessed fan artist filled an entire sketchbook page - messy, overlapping, full-body poses, tiny chibi doodles, exaggerated",
    "promptPreview": "Draw me as if an obsessed fan artist filled an entire sketchbook page - messy, overlapping, full-body poses, tiny chibi doodles, exaggerated expressions, and ra",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Ciri_ai",
    "tags": [
      "Characters & People",
      "Character",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-326",
    "locale": "zh",
    "title": "红蓝撞色高跟诱惑",
    "category": "角色与人物",
    "description": "[中文] { \"global_settings\": { \"resolution\": \"8K\", \"quality\": \"超高清晰度\", \"aspect_ratio\": \"2:3\", \"render_style\": \"AI编辑、高细节3D渲染\", \"lighting_quality",
    "promptPreview": "[中文] { \"global_settings\": { \"resolution\": \"8K\", \"quality\": \"超高清晰度\", \"aspect_ratio\": \"2:3\", \"render_style\": \"AI编辑、高细节3D渲染\", \"lighting_quality\": \"柔和影棚光与逼真阴影\", \"sh",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@meng_dagg695",
    "tags": [
      "Characters & People",
      "Realistic",
      "Character",
      "3D",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-25",
    "locale": "zh",
    "title": "综合应用场景图",
    "category": "角色与人物",
    "description": "create a minecraft skin inspired by {argument name=\"reference\" default=\"my look\"}",
    "promptPreview": "create a minecraft skin inspired by {argument name=\"reference\" default=\"my look\"}",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@nicdunz",
    "tags": [
      "Characters & People",
      "Characters",
      "Story"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-397",
    "locale": "zh",
    "title": "街舞角色设定参考图",
    "category": "角色与人物",
    "description": "角色设定图布局，聚焦于一位18岁的亚裔女性街舞舞者。包含4个大型、高细节度的全身动态舞姿（突出舞蹈动作，面部清晰）。侧边附一条清晰的多角度参考条，仅含3个精细头部特写（正面、侧面、3/4侧面）。最大限度减少文字元素，将像素空间优先用于面部细节刻画。背景为粗砺工业风，搭配写实光影效",
    "promptPreview": "角色设定图布局，聚焦于一位18岁的亚裔女性街舞舞者。包含4个大型、高细节度的全身动态舞姿（突出舞蹈动作，面部清晰）。侧边附一条清晰的多角度参考条，仅含3个精细头部特写（正面、侧面、3/4侧面）。最大限度减少文字元素，将像素空间优先用于面部细节刻画。背景为粗砺工业风，搭配写实光影效果",
    "sourceId": "freestylefly",
    "sourceLanguage": "ZH",
    "sourceLabel": "@ChangningL29508",
    "tags": [
      "Characters & People",
      "Realistic",
      "Character",
      "Creative"
    ],
    "featured": false,
    "needsReference": false,
    "chunkId": "chunk-5"
  },
  {
    "id": "freestylefly-502",
    "locale": "en",
    "title": "黑桃国王递归扑克牌",
    "category": "角色与人物",
    "description": "Use my uploaded face image as the primary identity reference. Preserve my exact facial identity with extremely high fidelity: identical faci",
    "promptPreview": "Use my uploaded face image as the primary identity reference. Preserve my exact facial identity with extremely high fidelity: identical facial structure, jawlin",
    "sourceId": "freestylefly",
    "sourceLanguage": "EN",
    "sourceLabel": "@Professor_134",
    "tags": [
      "Characters & People",
      "UI",
      "Poster",
      "Realistic",
      "Tech",
      "Commerce",
      "Fashion"
    ],
    "featured": false,
    "needsReference": true,
    "chunkId": "chunk-5"
  }
]
