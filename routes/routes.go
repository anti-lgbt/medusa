package routes

import (
	"github.com/anti-lgbt/medusa/controllers/admin"
	"github.com/anti-lgbt/medusa/controllers/identity"
	"github.com/anti-lgbt/medusa/controllers/public"
	"github.com/anti-lgbt/medusa/controllers/resource"
	"github.com/anti-lgbt/medusa/routes/middlewares"
	"github.com/gofiber/fiber/v2"

	fiberSwagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRouter() *fiber.App {
	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024, // this is the default limit of 4MB
	})
	app.Use(logger.New())

	app.Get("/api/v2/swagger/*", fiberSwagger.Handler)

	api_public := app.Group("/api/v2/public")
	{
		api_public.Get("time", public.GetTime)

		api_public_musics := api_public.Group("/musics")
		{
			api_public_musics.Get("/", public.GetMusics)
			api_public_musics.Get("/:id", public.GetMusic)
			api_public_musics.Get("/:id/image", public.GetMusicImage)
			api_public_musics.Get("/:id/audio", public.GetMusicAudio)
			api_public_musics.Get("/:id/comments", public.GetMusicComments)
		}

		api_public_albums := api_public.Group("/albums")
		{
			api_public_albums.Get("/:id", public.GetAlbum)
			api_public_albums.Get("/:id/image", public.GetAlbumImage)
			api_public_albums.Get("/:id/comments", public.GetAlbumComments)
		}

		api_public.Get("/users/:id/avatar", public.GetUserAvatar)
	}

	api_resource := app.Group("/api/v2/resource", middlewares.MustAuth)
	{
		api_resource_musics := api_resource.Group("/musics", middlewares.MustCollaborator)
		{
			api_resource_musics.Get("/", resource.GetMusics)
			api_resource_musics.Get("/:id", resource.GetMusic)
			api_resource_musics.Post("/", resource.CreateMusic)
			api_resource_musics.Put("/", resource.UpdateMusic)
			api_resource_musics.Delete("/:id", resource.DeleteMusic)
			api_resource_musics.Post("/:id/like", resource.LikeMusic)
			api_resource_musics.Post("/:id/unlike", resource.UnLikeMusic)
			api_resource_musics.Post("/:id/comment", resource.CommentMusic)
		}

		api_resource_albums := api_resource.Group("/albums")
		{
			api_resource_albums.Get("/", resource.GetAlbums)
			api_resource_albums.Get("/:id", resource.GetAlbum)
			api_resource_albums.Post("/", resource.CreateAlbum)
			api_resource_albums.Put("/:id", resource.UpdateAlbum)
			api_resource_albums.Delete("/:id", resource.DeleteAlbum)
			api_resource_albums.Post("/:id/like", resource.LikeAlbum)
			api_resource_albums.Post("/:id/unlike", resource.UnLikeAlbum)
			api_resource_albums.Post("/:id/comment", resource.CommentAlbum)
		}

		api_resource_users := api_resource.Group("/users")
		{
			api_resource_users.Get("/me", resource.GetUserProfile)
			api_resource_users.Put("/", resource.UpdateUser)
			api_resource_users.Put("/password", resource.UpdateUserPassword)
		}

		api_resource_comments := api_resource.Group("/comments")
		{
			api_resource_comments.Get("/:id/like", resource.LikeComment)
			api_resource_comments.Get("/:id/unlike", resource.UnLikeComment)
			api_resource_comments.Delete("/:id", resource.DeleteComment)

			api_resource_replys := api_resource_comments.Group("/reply")
			{
				api_resource_replys.Post("/", resource.CreateReply)
				api_resource_replys.Post("/:id/like", resource.LikeReply)
				api_resource_replys.Post("/:id/unlike", resource.UnLikeReply)
				api_resource_replys.Delete("/:id", resource.DeleteReply)
			}
		}
	}

	api_identity := app.Group("/api/v2/identity")
	{
		api_identity.Post("/session", middlewares.MustGuest, identity.Login)
		api_identity.Delete("/session", middlewares.MustAuth, identity.Logout)

		api_identity.Post("/users", middlewares.MustGuest, identity.Register)
		api_identity.Post("/users/generate_code", identity.ReSendEmailCode)
		api_identity.Post("/users/confirm_code", identity.VerifyEmail)

		api_identity.Post("/password/generate_code", middlewares.MustGuest, identity.GenerateCodeResetPassword)
		api_identity.Post("/password/check_code", middlewares.MustGuest, identity.CheckCodeResetPassword)
		api_identity.Post("/password/reset_password", middlewares.MustGuest, identity.ResetPassword)
	}

	api_admin := app.Group("/api/v2/admin", middlewares.MustAuth, middlewares.MustAdmin)
	{
		api_admin_users := api_admin.Group("/users")
		{
			api_admin_users.Get("/", admin.GetUsers)
			api_admin_users.Get("/:uid", admin.GetUser)
			api_admin_users.Put("/:uid", admin.UpdateUser)
			api_admin_users.Delete("/:uid", admin.DeleteUser)
		}
	}

	return app
}
