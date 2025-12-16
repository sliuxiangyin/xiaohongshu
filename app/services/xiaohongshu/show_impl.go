package xiaohongshu

type ShowImpl[T any] interface {
	Show() (T, error)
}
