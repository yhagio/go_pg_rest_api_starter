package models

type postValidator struct {
	PostDB
}

func (pv *postValidator) Create(post *Post) error {
	err := postValidationFuncs(post,
		pv.userIDRequired,
		pv.titleRequired,
		pv.descRequired)

	if err != nil {
		return err
	}

	return pv.PostDB.Create(post)
}

func (pv *postValidator) Update(post *Post) error {
	err := postValidationFuncs(post,
		pv.userIDRequired,
		pv.titleRequired,
		pv.descRequired)

	if err != nil {
		return err
	}

	return pv.PostDB.Update(post)
}

func (pv *postValidator) Delete(id uint) error {
	var post Post
	post.ID = id

	err := postValidationFuncs(&post, pv.validId)

	if err != nil {
		return err
	}

	return pv.PostDB.Delete(post.ID)
}

///////////////////////////////////////////////////////////
// Private functions
///////////////////////////////////////////////////////////

func (pv *postValidator) userIDRequired(post *Post) error {
	if post.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (pv *postValidator) titleRequired(post *Post) error {
	if post.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

func (pv *postValidator) descRequired(post *Post) error {
	if post.Description == "" {
		return ErrDescRequired
	}
	return nil
}

func (pv *postValidator) validId(post *Post) error {
	if post.ID <= 0 {
		return ErrInvalidID
	}
	return nil
}

///////////////////////////////////////////////////////////
// Reusable validation functions helper
///////////////////////////////////////////////////////////

type postValidationFunc func(*Post) error

func postValidationFuncs(post *Post, funcs ...postValidationFunc) error {
	for _, fn := range funcs {
		err := fn(post)
		if err != nil {
			return err
		}
	}
	return nil
}
