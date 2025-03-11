package secure

import "iis_server/apiq"

var instagramUsers = make(map[string]apiq.InstagramUsername)

func GetAllUsers() []apiq.InstagramUsername {
	users := []apiq.InstagramUsername{}
	for _, user := range instagramUsers {
		users = append(users, user)
	}
	return users
}

func GetUserByID(userID string) (apiq.InstagramUsername, bool) {
	user, exists := instagramUsers[userID]
	return user, exists
}

func CreateUser(user apiq.InstagramUsername) {
	instagramUsers[user.UserId] = user
}

func UpdateUser(userID string, updatedUser apiq.InstagramUsername) bool {
	_, exists := instagramUsers[userID]
	if exists {
		instagramUsers[userID] = updatedUser
	}
	return exists
}

func DeleteUser(userID string) bool {
	_, exists := instagramUsers[userID]
	if exists {
		delete(instagramUsers, userID)
	}
	return exists
}
