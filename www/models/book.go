package models

type D map[string]interface{}

var (
	Tags = map[string][]string{
		"文学": []string{"小说", "外国文学", "文学", "随笔", "中国文学", "经典", "散文", "日本文学", "村上春树", "童话", "诗歌", "王小波", "杂文", "张爱玲", "儿童文学", "余华", "古典文学", "名著", "钱钟书", "当代文学", "鲁迅", "外国名著", "诗词", "茨威格", "杜拉斯", "米兰·昆德拉", "港台"},
		"流行": []string{"漫画", "绘本", "推理", "青春", "言情", "科幻", "韩寒", "武侠", "悬疑", "耽美", "亦舒", "东野圭吾", "日本漫画", "奇幻", "安妮宝贝", "三毛", "郭敬明", "网络小说", "穿越", "金庸", "几米", "轻小说", "推理小说", "阿加莎·克里斯蒂", "张小娴", "幾米", "魔幻", "青春文学", "高木直子", "J.K.罗琳", "沧月", "落落", "张悦然", "古龙", "科幻小说", "蔡康永"},
		"文化": []string{"历史", "心理学", "哲学", "传记", "文化", "社会学", "设计", "艺术", "政治", "社会", "建筑", "宗教", "电影", "数学", "政治学", "回忆录", "思想", "国学", "中国历史", "音乐", "人文", "戏剧", "人物传记", "绘画", "艺术史", "佛教", "军事", "西方哲学", "二战", "自由主义", "近代史", "考古", "美术"},
		"生活": []string{"爱情", "旅行", "生活", "励志", "成长", "摄影", "心理", "女性", "职场", "美食", "游记", "教育", "灵修", "情感", "健康", "手工", "养生", "两性", "家居", "人际关系", "自助游"},
		"经管": []string{"经济学", "管理", "经济", "金融", "商业", "投资", "营销", "理财", "创业", "广告", "股票", "企业史", "策划"},
		"科技": []string{"科普", "互联网", "编程", "科学", "交互设计", "用户体验", "算法", "web", "科技", "UE", "UCD", "通信", "交互", "神经网络", "程序"},
	}
)

func GetBookCatalog() (catalog []string) {
	for k, _ := range Tags {
		catalog = append(catalog, k)
	}
	return
}

type Book struct {
	Catalog      string
	Tag          string
	Id           string
	Isbn10       string
	Isbn13       string
	Title        string
	Origin_title string
	Alt_title    string
	Subtitle     string
	Url          string
	Alt          string
	Image        string
	Images       D
	Author       []string
	Translator   []string
	Publisher    string
	Pubdate      string
	Rating       D
	Tags         string
	Binding      string
	Price        string
	Pages        string
	Author_intro string
	Summary      string
	In           []string
	Out          []string
}
