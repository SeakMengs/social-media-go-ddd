"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { PostCard } from "@/components/posts/post-card";
import { Loader2, Bookmark } from "lucide-react";
import { AuthUserResult } from "@/auth";
import { AggregatePost, PostType } from "@/types/model";
import { getMyFavoritePosts } from "@/api/action";
import { toast } from "sonner";
import { getPostKey } from "@/utils/post";

type FavoritesProps = {
  auth: AuthUserResult;
};

export default function FavoritesPage({ auth }: FavoritesProps) {
  const user = auth && auth.user;
  const [favoritePosts, setFavoritePosts] = useState<AggregatePost[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchFavoritePosts = async () => {
      if (!user) return;

      try {
        const response = await getMyFavoritePosts();
        console.log("Fetched favorite posts:", response);
        if (response.success) {
          setFavoritePosts(response.data.posts);
        } else {
          toast.error(response.error || "Failed to fetch favorite posts");
        }
      } catch (error) {
      if (error instanceof Error) {
        toast.error(error.message || "An unexpected error occurred");
      }
      } finally {
        setIsLoading(false);
      }
    };

    fetchFavoritePosts();
  }, [user]);

  const handlePostUpdate = (updatedPost: AggregatePost) => {
    const updatedTargetId =
      updatedPost.type === PostType.REPOST
        ? updatedPost.repost?.postId || updatedPost.id
        : updatedPost.id;

    setFavoritePosts((prevPosts) => {
      // If post is unfavorited, remove all posts that reference the same original content
      if (!updatedPost.favorited) {
        return prevPosts.filter((post) => {
          const postTargetId = post.type === PostType.REPOST && post.repost 
            ? post.repost.postId 
            : post.id;
          
          return postTargetId !== updatedTargetId;
        });
      }
      
      // Otherwise update interactions for all posts with the same original content
      return prevPosts.map((post) => {
        const postTargetId =
          post.type === PostType.REPOST
            ? post.repost?.postId || post.id
            : post.id;

        if (postTargetId === updatedTargetId) {
          return {
            ...post,
            liked: updatedPost.liked,
            favorited: updatedPost.favorited,
            reposted: updatedPost.reposted,
            likeCount: updatedPost.likeCount,
            favoriteCount: updatedPost.favoriteCount,
            repostCount: updatedPost.repostCount,
          };
        }
        
        return post;
      });
    });
  };

  const handleUnrepost = (unrepostedPost: AggregatePost) => {
    // For favorites, we don't need to remove reposts, just sync the state
    // The handlePostUpdate will handle the interaction syncing
  };

  const handleRepostCreated = (newRepost: AggregatePost) => {
    // For favorites, we don't need to do anything special for new reposts
    // The handlePostUpdate will handle the interaction syncing
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background">
        <div className="container max-w-2xl mx-auto py-8 px-4">
          <div className="mb-8">
            <h1 className="text-3xl font-heading font-bold text-primary">
              Your Favorites
            </h1>
            <p className="text-muted-foreground mt-2">
              Posts you've saved for later
            </p>
          </div>
          <div className="flex items-center justify-center py-8">
            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <div className="container max-w-2xl mx-auto py-8 px-4">
        <div className="mb-8">
          <h1 className="text-3xl font-heading font-bold text-primary">
            Your Favorites
          </h1>
          <p className="text-muted-foreground mt-2">
            Posts you've saved for later
          </p>
        </div>

        {favoritePosts.length > 0 ? (
          <div className="space-y-4">
            {favoritePosts.map((post, i) => (
              <PostCard
                key={getPostKey(post, i)}
                post={post}
                onPostUpdate={handlePostUpdate}
                onUnrepost={handleUnrepost}
                onRepostCreated={handleRepostCreated}
                auth={auth}
              />
            ))}
          </div>
        ) : (
          <div className="text-center py-12">
            <Bookmark className="h-16 w-16 mx-auto mb-4 text-muted-foreground opacity-50" />
            <h2 className="text-xl font-heading font-semibold text-foreground mb-2">
              No favorites yet
            </h2>
            <p className="text-muted-foreground">
              When you favorite posts, they'll appear here for easy access
              later.
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
