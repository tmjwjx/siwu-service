package globals

// Home :图片所属单位，即属于文章图片还是用户图片或是资源表图片
type Home string

const (
	UserHome          Home = "user"
	ArticleHome       Home = "article"
	TagHome           Home = "tag"
	CategoryHome      Home = "category"
	AdvertisementHome Home = "advertisement"
	CommentHome       Home = "comment"
)
