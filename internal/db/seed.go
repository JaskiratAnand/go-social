package db

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/JaskiratAnand/go-social/internal/store"
	"github.com/google/uuid"
)

var (
	NumberOfUsers    = 200
	NumberOfPosts    = 250
	NumberOfComments = 1500
	NumberOfFollows  = 500
)

func Seed(store *store.Queries) error {

	ctx := context.Background()

	log.Println("generating users...")
	users := generateUsers(NumberOfUsers)
	userIDs := make([]uuid.UUID, NumberOfUsers)
	for i, user := range users {
		var err error
		userIDs[i], err = store.CreateUser(ctx, *user)
		if err != nil {
			log.Println("Error creating user while seeding data")
			return err
		}
	}

	log.Println("generating posts...")
	posts := generatePosts(NumberOfPosts, userIDs)
	postIDs := make([]uuid.UUID, NumberOfPosts)
	for i, post := range posts {
		createPost, err := store.CreatePost(ctx, *post)
		if err != nil {
			log.Println("Error creating posts while seeding data")
			return err
		}
		postIDs[i] = createPost.ID
	}

	log.Println("generating comments...")
	comments := generateComments(NumberOfComments, userIDs, postIDs)
	for _, comment := range comments {
		_, err := store.CreateComment(ctx, *comment)
		if err != nil {
			log.Println("Error creating comments while seeding data")
			return err
		}
	}

	log.Println("generating follows...")
	follows := generateFollows(NumberOfFollows, userIDs)
	for _, follow := range follows {
		err := store.FollowUser(ctx, *follow)
		if err != nil {
			log.Println("Error creating follows while seeding data")
			return err
		}
	}

	return nil
}

func generateUsers(num int) []*store.CreateUserParams {

	names := []string{
		"Aarav", "Emma", "Liam", "Sophia", "Noah", "Olivia", "Ethan", "Ava", "Mason", "Isabella",
		"Lucas", "Mia", "Elijah", "Amelia", "Logan", "Harper", "James", "Charlotte", "Aiden", "Ella",
		"Jackson", "Lily", "Alexander", "Aria", "Benjamin", "Chloe", "Sebastian", "Zoey", "William", "Grace",
		"Henry", "Hannah", "Gabriel", "Ellie", "Matthew", "Scarlett", "Daniel", "Victoria", "Michael", "Layla",
		"Samuel", "Nora", "David", "Hazel", "Joseph", "Aurora", "Carter", "Riley", "Owen", "Violet",
	}

	users := make([]*store.CreateUserParams, num)
	for i := 0; i < num; i++ {
		name := names[rand.IntN(len(names))]
		username := fmt.Sprintf("%s_%v", name, i+1)

		users[i] = &store.CreateUserParams{
			Username: username,
			Email:    fmt.Sprintf("%v@example.com", username),
			Password: []byte(""),
		}
	}
	return users
}

func generatePosts(num int, userIDs []uuid.UUID) []*store.CreatePostParams {

	titles := []string{
		"The Power of Habit", "Embracing Minimalism", "Healthy Eating Tips",
		"Travel on a Budget", "Mindfulness Meditation", "Boost Your Productivity",
		"Home Office Setup", "Digital Detox", "Gardening Basics",
		"DIY Home Projects", "Yoga for Beginners", "Sustainable Living",
		"Mastering Time Management", "Exploring Nature", "Simple Cooking Recipes",
		"Fitness at Home", "Personal Finance Tips", "Creative Writing",
		"Mental Health Awareness", "Learning New Skills",
	}
	contents := []string{
		"In this post, we'll explore how to develop good habits that stick and transform your life.",
		"Discover the benefits of a minimalist lifestyle and how to declutter your home and mind.",
		"Learn practical tips for eating healthy on a budget without sacrificing flavor.",
		"Traveling doesn't have to be expensive. Here are some tips for seeing the world on a budget.",
		"Mindfulness meditation can reduce stress and improve your mental well-being. Here's how to get started.",
		"Increase your productivity with these simple and effective strategies.",
		"Set up the perfect home office to boost your work-from-home efficiency and comfort.",
		"A digital detox can help you reconnect with the real world and improve your mental health.",
		"Start your gardening journey with these basic tips for beginners.",
		"Transform your home with these fun and easy DIY projects.",
		"Yoga is a great way to stay fit and flexible. Here are some beginner-friendly poses to try.",
		"Sustainable living is good for you and the planet. Learn how to make eco-friendly choices.",
		"Master time management with these tips and get more done in less time.",
		"Nature has so much to offer. Discover the benefits of spending time outdoors.",
		"Whip up delicious meals with these simple and quick cooking recipes.",
		"Stay fit without leaving home with these effective at-home workout routines.",
		"Take control of your finances with these practical personal finance tips.",
		"Unleash your creativity with these inspiring writing prompts and exercises.",
		"Mental health is just as important as physical health. Learn how to take care of your mind.",
		"Learning new skills can be fun and rewarding. Here are some ideas to get you started.",
	}
	tags := []string{
		"Self Improvement", "Minimalism", "Health", "Travel", "Mindfulness",
		"Productivity", "Home Office", "Digital Detox", "Gardening", "DIY",
		"Yoga", "Sustainability", "Time Management", "Nature", "Cooking",
		"Fitness", "Personal Finance", "Writing", "Mental Health", "Learning",
	}

	posts := make([]*store.CreatePostParams, num)
	for i := 0; i < num; i++ {
		posts[i] = &store.CreatePostParams{
			Title:   fmt.Sprintf("title_%v %v", i, titles[rand.IntN(len(titles))]),
			Content: contents[rand.IntN(len(contents))],
			UserID:  userIDs[rand.IntN(NumberOfUsers)],
			Tags:    tags[:rand.IntN(len(tags))],
		}
	}
	return posts
}

func generateComments(num int, userIDs, postIDs []uuid.UUID) []*store.CreateCommentParams {

	contents := []string{
		"Great post! Thanks for sharing.",
		"I completely agree with your thoughts.",
		"Thanks for the tips, very helpful.",
		"Interesting perspective, I hadn't considered that.",
		"Thanks for sharing your experience.",
		"Well written, I enjoyed reading this.",
		"This is very insightful, thanks for posting.",
		"Great advice, I'll definitely try that.",
		"I love this, very inspirational.",
		"Thanks for the information, very useful.",
	}

	comments := make([]*store.CreateCommentParams, num)
	for i := 0; i < num; i++ {
		comments[i] = &store.CreateCommentParams{
			PostID:  postIDs[rand.IntN(NumberOfPosts)],
			UserID:  userIDs[rand.IntN(NumberOfUsers)],
			Content: contents[rand.IntN(len(contents))],
		}
	}
	return comments
}

func generateFollows(num int, userIDs []uuid.UUID) []*store.FollowUserParams {
	existingPairs := make(map[string]struct{})
	follows := make([]*store.FollowUserParams, 0, num)

	for len(follows) < num {
		userID := userIDs[rand.IntN(len(userIDs))]
		followID := userIDs[rand.IntN(len(userIDs))]
		if userID == followID {
			continue
		}

		pairKey := fmt.Sprintf("%s:%s", userID, followID)
		if _, exists := existingPairs[pairKey]; exists {
			continue
		}

		existingPairs[pairKey] = struct{}{}
		follows = append(follows, &store.FollowUserParams{
			UserID:   userID,
			FollowID: followID,
		})
	}
	return follows
}
