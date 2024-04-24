package utils

import "serve/model"

func IsOutsideArea(pos model.Position, border model.Border) bool {
	return pos.X < border.W+border.X && pos.X > border.X && pos.Y < border.H+border.Y && pos.Y > border.Y
}
