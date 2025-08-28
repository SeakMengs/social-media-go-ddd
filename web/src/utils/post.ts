import { AggregatePost, PostType } from "@/types/model";

export function getPostKey(post: AggregatePost, i: number | null = null) {
  // Always use the post's own unique ID as the primary key
  // Each post (including reposts) has its own unique ID
  return post.type === PostType.REPOST && post.repost
    ? post.repost.id
    : post.id + (i !== null ? `-${i}` : "");
}
export function getPostId(post: AggregatePost) {
  return post.type === PostType.REPOST && post.repost
    ? post.repost.postId
    : post.id;
}
