import { z } from "zod";

export const BaseModelSchema = z.object({
    id: z.string(),
    createdAt: z.coerce.date(),
    updatedAt: z.coerce.date(),
});
export type BaseModel = z.infer<typeof BaseModelSchema>;

export const UserSchema = BaseModelSchema.extend({
    username: z.string(),
    email: z.email(),
});
export type User = z.infer<typeof UserSchema>;

export const PostSchema = BaseModelSchema.extend({
    userId: z.string(),
    content: z.string(),
});
export type Post = z.infer<typeof PostSchema>;

export const LikeSchema = BaseModelSchema.extend({
    userId: z.string(),
    postId: z.string(),
});
export type Like = z.infer<typeof LikeSchema>;

export const FavoriteSchema = BaseModelSchema.extend({
    userId: z.string(),
    postId: z.string(),
});
export type Favorite = z.infer<typeof FavoriteSchema>;

export const RepostSchema = BaseModelSchema.extend({
    userId: z.string(),
    postId: z.string(),
    comment: z.string(),
});
export type Repost = z.infer<typeof RepostSchema>;

export const SessionSchema = BaseModelSchema.extend({
    userId: z.string(),
    expireAt: z.coerce.date(),
});
export type Session = z.infer<typeof SessionSchema>;

export const FollowSchema = BaseModelSchema.extend({
    followerId: z.string(),
    followeeId: z.string(),
});
export type Follow = z.infer<typeof FollowSchema>;
