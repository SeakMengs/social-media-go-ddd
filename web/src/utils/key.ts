import { AggregatePost, PostType } from "@/types/model";

export function getPostKey(post: AggregatePost) {
  return post.type === PostType.REPOST && post.repost ? post.repost.id : post.id;
}