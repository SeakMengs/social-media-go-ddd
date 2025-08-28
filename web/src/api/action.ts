"use server";

import {
  AggregatePost,
  AggregatePostSchema,
  Post,
  User,
  UserSchema,
  Session,
  SessionSchema,
  Repost,
  AggregateUser,
  AggregateUserSchema,
} from "@/types/model";
import { apiWithAuth, api } from "@/utils/axios";
import { createScopedLogger } from "@/utils/logger";
import { ResponseJson, responseSomethingWentWrong } from "@/utils/response";
import z from "zod";

const logger = createScopedLogger("api:user");

// -- Auth Endpoints ---
export type RegisterUserRequest = {
  email: string;
  password: string;
  username: string;
};
export type RegisterUserResponse = {
  user: User;
};
export async function registerUser(
  data: RegisterUserRequest
): Promise<ResponseJson<RegisterUserResponse>> {
  try {
    const response = await api.post<ResponseJson<RegisterUserResponse>>(
      "/api/v1/auth/register",
      data
    );
    return response.data;
  } catch (error) {
    logger.error("Error registering user", error);
    return responseSomethingWentWrong("Error registering user (catch)");
  }
}

// --- User Endpoints ---

// Get Current user posts
export type GetMyPostsResponse = AggregatePost[];
export async function getMyPosts(): Promise<ResponseJson<GetMyPostsResponse>> {
  try {
    const response = await apiWithAuth.get<ResponseJson<GetMyPostsResponse>>(
      "/api/v1/users/me/posts"
    );
    if (!response.data.success) {
      logger.debug("Get my posts failed", response.data.error);
      return response.data;
    }
    const validate = AggregatePostSchema.array().safeParse(response.data.data);
    if (!validate.success) {
      logger.debug("Get my posts failed", validate.error);

      return response.data;
    }
    return response.data;
  } catch (error) {
    logger.error("Error getting my posts", error);
    return responseSomethingWentWrong("Error getting my posts (catch)");
  }
}

// Get Current User Profile
export type GetCurrentUserProfileResponse = {
  session: Session;
  user: AggregateUser;
};
export async function getCurrentUserProfile(): Promise<
  ResponseJson<GetCurrentUserProfileResponse>
> {
  try {
    const response = await apiWithAuth.get<
      ResponseJson<GetCurrentUserProfileResponse>
    >("/api/v1/users/me");
    if (!response.data.success) {
      logger.debug("Get current user profile failed", response.data.error);
      return response.data;
    }
    const schema = z.object({
      session: SessionSchema,
      user: AggregateUserSchema,
    });
    const validate = schema.safeParse(response.data.data);
    if (!validate.success) {
      logger.debug("Get current user profile failed", validate.error);
      return response.data;
    }
    return response.data;
  } catch (error) {
    logger.error("Error getting current user profile", error);
    return responseSomethingWentWrong(
      "Error getting current user profile (catch)"
    );
  }
}

// Get User by ID
export type GetUserByIdResponse = { user: AggregateUser };
export async function getUserById(
  id: string
): Promise<ResponseJson<GetUserByIdResponse>> {
  try {
    const response = await apiWithAuth.get<ResponseJson<GetUserByIdResponse>>(
      `/api/v1/public/users/${id}`
    );
    if (!response.data.success) {
      logger.debug("Get user by id failed", response.data.error);
      return response.data;
    }
    const schema = z.object({ user: AggregateUserSchema });
    const validate = schema.safeParse(response.data.data);
    if (!validate.success) {
      logger.debug("Get user by id failed", validate.error);
      return response.data;
    }
    return response.data;
  } catch (error) {
    logger.error("Error getting user by id", error);
    return responseSomethingWentWrong("Error getting user by id (catch)");
  }
}

export type SearchUsersByNameResponse = { users: AggregateUser[] };
export async function searchUsersByName(
  name: string
): Promise<ResponseJson<SearchUsersByNameResponse>> {
  try {
    const response = await apiWithAuth.get<
      ResponseJson<SearchUsersByNameResponse>
    >(`/api/v1/public/users/name/${name}`);
    if (!response.data.success) {
      logger.debug("Search user by name failed", response.data.error);
      return response.data;
    }
    const schema = z.object({ users: AggregateUserSchema.array().default([]) });
    const validate = schema.safeParse(response.data.data);
    if (!validate.success) {
      logger.debug("Search user by name failed", validate.error);
      return response.data;
    }
    return response.data;
  } catch (error) {
    logger.error("Error getting user by name", error);
    return responseSomethingWentWrong("Error getting user by name (catch)");
  }
}

// Follow User
export async function followUser(id: string): Promise<ResponseJson<{}>> {
  try {
    const response = await apiWithAuth.post<ResponseJson<{}>>(
      `/api/v1/users/${id}/follow`
    );
    return response.data;
  } catch (error) {
    logger.error("Error following user", error);
    return responseSomethingWentWrong("Error following user (catch)");
  }
}

// Unfollow User
export async function unfollowUser(id: string): Promise<ResponseJson<{}>> {
  try {
    const response = await apiWithAuth.delete<ResponseJson<{}>>(
      `/api/v1/users/${id}/follow`
    );
    return response.data;
  } catch (error) {
    logger.error("Error unfollowing user", error);
    return responseSomethingWentWrong("Error unfollowing user (catch)");
  }
}

// Get My Reposts
export type GetMyRepostsResponse = { posts: AggregatePost[] };
export async function getMyReposts(): Promise<
  ResponseJson<GetMyRepostsResponse>
> {
  try {
    const response = await apiWithAuth.get<ResponseJson<GetMyRepostsResponse>>(
      "/api/v1/users/me/reposts"
    );
    if (!response.data.success) {
      logger.debug("Get my reposts failed", response.data.error);
      return response.data;
    }
    const validate = z
      .object({
        posts: AggregatePostSchema.array().default([]),
      })
      .safeParse(response.data.data);
    if (!validate.success) {
      logger.debug("Get my reposts failed", validate.error);
      return response.data;
    }
    return response.data;
  } catch (error) {
    logger.error("Error getting my reposts", error);
    return responseSomethingWentWrong("Error getting my reposts (catch)");
  }
}

// Get User Posts
export type GetUserPostsResponse = { posts: AggregatePost[] };
export async function getUserPosts(
  id: string
): Promise<ResponseJson<GetUserPostsResponse>> {
  try {
    const response = await api.get<ResponseJson<GetUserPostsResponse>>(
      `/api/v1/public/users/${id}/posts`
    );
    if (!response.data.success) {
      logger.debug("Get user posts failed", response.data.error);
      return response.data;
    }
    const validate = z
      .object({
        posts: AggregatePostSchema.array().default([]),
      })
      .safeParse(response.data.data);
    if (!validate.success) {
      logger.debug("Get user posts failed", validate.error);
      return response.data;
    }
    return response.data;
  } catch (error) {
    logger.error("Error getting user posts", error);
    return responseSomethingWentWrong("Error getting user posts (catch)");
  }
}
// Get User Reposts
export type GetUserRepostsResponse = { reposts: AggregatePost[] };
export async function getUserReposts(
  id: string
): Promise<ResponseJson<GetUserRepostsResponse>> {
  try {
    const response = await api.get<ResponseJson<GetUserRepostsResponse>>(
      `/api/v1/public/users/${id}/reposts`
    );
    if (!response.data.success) {
      logger.debug("Get user reposts failed", response.data.error);
      return response.data;
    }
    const validate = z
      .object({
        reposts: AggregatePostSchema.array().default([]),
      })
      .safeParse(response.data.data);
    if (!validate.success) {
      logger.debug("Get user reposts failed", validate.error);
      return response.data;
    }
    return response.data;
  } catch (error) {
    logger.error("Error getting user reposts", error);
    return responseSomethingWentWrong("Error getting user reposts (catch)");
  }
}
// Get my FavoritePosts
export type GetMyFavoritePostsResponse = { posts: AggregatePost[] };
export async function getMyFavoritePosts(): Promise<
  ResponseJson<GetMyFavoritePostsResponse>
> {
  try {
    const response = await apiWithAuth.get<
      ResponseJson<GetMyFavoritePostsResponse>
    >(`/api/v1/users/me/posts/favorites`);
    if (!response.data.success) {
      logger.debug("Get my favorite posts failed", response.data.error);
      return response.data;
    }
    const validate = z
      .object({
        posts: AggregatePostSchema.array().default([]),
      })
      .safeParse(response.data.data);
    if (!validate.success) {
      logger.debug("Get my favorite posts failed", validate.error);
      return response.data;
    }
    return response.data;
  } catch (error) {
    logger.error("Error getting my favorite posts", error);
    return responseSomethingWentWrong(
      "Error getting my favorite posts (catch)"
    );
  }
}

// Get My Feed
export type FeedPagination = {
  hasMore: any;
  page: number;
  pageSize: number;
  total: number;
};
export type GetMyFeedResponse = {
  feed: AggregatePost[];
  pagination: FeedPagination;
};
export async function getMyFeed(params?: {
  page?: number;
  pageSize?: number;
}): Promise<ResponseJson<GetMyFeedResponse>> {
  try {
    const response = await apiWithAuth.get<ResponseJson<GetMyFeedResponse>>(
      "/api/v1/users/me/feed",
      { params }
    );
    if (!response.data.success) {
      logger.debug("Get my feed failed", response.data.error);
      return response.data;
    }
    const validate = z
      .object({
        feed: AggregatePostSchema.array(),
        pagination: z.object({
          page: z.number(),
          pageSize: z.number(),
          total: z.number(),
        }),
      })
      .safeParse(response.data.data);
    if (!validate.success) {
      logger.debug("Get my feed failed", validate.error);
      return response.data;
    }
    return response.data;
  } catch (error) {
    logger.error("Error getting my feed", error);
    return responseSomethingWentWrong("Error getting my feed (catch)");
  }
}

// --- Post Endpoints ---

// Create Post
export type CreatePostRequest = { content: string };
export type CreatePostResponse = { post: AggregatePost };
export async function createPost(
  data: CreatePostRequest
): Promise<ResponseJson<CreatePostResponse>> {
  try {
    const response = await apiWithAuth.post<ResponseJson<CreatePostResponse>>(
      "/api/v1/posts",
      data
    );
    return response.data;
  } catch (error) {
    logger.error("Error creating post", error);
    return responseSomethingWentWrong("Error creating post (catch)");
  }
}

// Get Post by ID
export type GetPostByIdResponse = { post: AggregatePost };
export async function getPostById(
  id: string
): Promise<ResponseJson<GetPostByIdResponse>> {
  try {
    const response = await api.get<ResponseJson<GetPostByIdResponse>>(
      `/api/v1/posts/${id}`
    );
    if (!response.data.success) {
      logger.debug("Get post by id failed", response.data.error);
      return response.data;
    }
    const validate = z
      .object({ post: AggregatePostSchema })
      .safeParse(response.data.data);
    if (!validate.success) {
      logger.debug("Get post by id failed", validate.error);
      return response.data;
    }
    return response.data;
  } catch (error) {
    logger.error("Error getting post by id", error);
    return responseSomethingWentWrong("Error getting post by id (catch)");
  }
}

// Update Post
export type UpdatePostRequest = { content: string };
export type UpdatePostResponse = { post: AggregatePost };
export async function updatePost(
  id: string,
  data: UpdatePostRequest
): Promise<ResponseJson<UpdatePostResponse>> {
  try {
    const response = await apiWithAuth.put<ResponseJson<UpdatePostResponse>>(
      `/api/v1/posts/${id}`,
      data
    );
    return response.data;
  } catch (error) {
    logger.error("Error updating post", error);
    return responseSomethingWentWrong("Error updating post (catch)");
  }
}

// Delete Post
export async function deletePost(id: string): Promise<ResponseJson<{}>> {
  try {
    const response = await apiWithAuth.delete<ResponseJson<{}>>(
      `/api/v1/posts/${id}`
    );
    return response.data;
  } catch (error) {
    logger.error("Error deleting post", error);
    return responseSomethingWentWrong("Error deleting post (catch)");
  }
}

// Like Post
export async function likePost(id: string): Promise<ResponseJson<{}>> {
  try {
    const response = await apiWithAuth.post<ResponseJson<{}>>(
      `/api/v1/posts/${id}/like`
    );
    return response.data;
  } catch (error) {
    logger.error("Error liking post", error);
    return responseSomethingWentWrong("Error liking post (catch)");
  }
}

// Unlike Post
export async function unlikePost(id: string): Promise<ResponseJson<{}>> {
  try {
    const response = await apiWithAuth.delete<ResponseJson<{}>>(
      `/api/v1/posts/${id}/like`
    );
    return response.data;
  } catch (error) {
    logger.error("Error unliking post", error);
    return responseSomethingWentWrong("Error unliking post (catch)");
  }
}

// Favorite Post
export async function favoritePost(id: string): Promise<ResponseJson<{}>> {
  try {
    const response = await apiWithAuth.post<ResponseJson<{}>>(
      `/api/v1/posts/${id}/favorite`
    );
    return response.data;
  } catch (error) {
    logger.error("Error favoriting post", error);
    return responseSomethingWentWrong("Error favoriting post (catch)");
  }
}

// Unfavorite Post
export async function unfavoritePost(id: string): Promise<ResponseJson<{}>> {
  try {
    const response = await apiWithAuth.delete<ResponseJson<{}>>(
      `/api/v1/posts/${id}/favorite`
    );
    return response.data;
  } catch (error) {
    logger.error("Error unfavoriting post", error);
    return responseSomethingWentWrong("Error unfavoriting post (catch)");
  }
}

// Repost
export type RepostRequest = { comment: string };
export type RepostResponse = { repost: Repost };
export async function repost(
  id: string,
  data: RepostRequest
): Promise<ResponseJson<RepostResponse>> {
  try {
    const response = await apiWithAuth.post<ResponseJson<RepostResponse>>(
      `/api/v1/posts/${id}/repost`,
      data
    );
    return response.data;
  } catch (error) {
    logger.error("Error reposting", error);
    return responseSomethingWentWrong("Error reposting (catch)");
  }
}

// Unrepost
export async function unrepost(id: string): Promise<ResponseJson<{}>> {
  try {
    const response = await apiWithAuth.delete<ResponseJson<{}>>(
      `/api/v1/reposts/${id}`
    );
    return response.data;
  } catch (error) {
    logger.error("Error unreposting", error);
    return responseSomethingWentWrong("Error unreposting (catch)");
  }
}
