"use client";

import { useState, useEffect } from "react";
import { PostCard } from "@/components/posts/post-card";
import { CreatePost } from "@/components/posts/create-post";
import { EditPostDialog } from "@/components/posts/edit-post-dialog";
import { Button } from "@/components/ui/button";
import { Loader2, RefreshCw } from "lucide-react";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { AggregatePost, PostType } from "@/types/model";
import { AuthUser } from "@/auth";
import { deletePost, getMyFeed } from "@/api/action";
import { toast } from "sonner";
import { getPostKey } from "@/utils/post";

type PersonalizedFeedProps = {
  auth: AuthUser;
};

export function PersonalizedFeed({ auth }: PersonalizedFeedProps) {
  const [posts, setPosts] = useState<AggregatePost[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [editingPost, setEditingPost] = useState<AggregatePost | null>(null);
  const [deletingPostId, setDeletingPostId] = useState<string | null>(null);
  const [pagination, setPagination] = useState({
    page: 1,
    pageSize: 10,
    total: 0,
    hasMore: false,
  });
  const { user } = auth;

  const fetchFeed = async (page = 1, append = false) => {
    if (!user) return;

    if (page === 1 && !append) {
      setIsLoading(true);
    } else {
      setIsLoadingMore(true);
    }

    try {
      const response = await getMyFeed({
        page,
        pageSize: pagination.pageSize,
      });
      if (response.success) {
        if (response.data.feed) {
          if (append) {
            setPosts((prev) => [...prev, ...response.data!.feed!]);
          } else {
            setPosts(response.data.feed);
          }
        }
        if (response.data.pagination) {
            setPagination({
            ...response.data.pagination,
            hasMore:
              response.data.pagination.page * response.data.pagination.pageSize <
              response.data.pagination.total,
            });
        }
      } else {
        toast.error(response.error || "Failed to load feed");
      }
    } catch (error) {
      toast.error("Failed to load feed");
    } finally {
      setIsLoading(false);
      setIsLoadingMore(false);
      setIsRefreshing(false);
    }
  };

  useEffect(() => {
    fetchFeed();
  }, [user]);

  const handlePostCreated = () => {
    handleRefresh();
  };

  const handlePostUpdated = () => {
    fetchFeed();
  };

  const handleRefresh = async () => {
    setIsRefreshing(true);
    await fetchFeed(1, false);
  };

  const handleLoadMore = () => {
    if (pagination.hasMore && !isLoadingMore) {
      fetchFeed(pagination.page + 1, true);
    }
  };

  const handleEdit = (post: AggregatePost) => {
    setEditingPost(post);
  };

  const handleDelete = (postId: string) => {
    setDeletingPostId(postId);
  };

  const confirmDelete = async () => {
    if (!deletingPostId || !user) return;

    try {
      const response = await deletePost(deletingPostId);
      if (response.success) {
        setPosts(posts.filter((p) => p.id !== deletingPostId));
        toast.success("Post deleted successfully");
      } else {
        toast.error(response.error);
      }
    } catch (error) {
      toast.error("Failed to delete post");
    } finally {
      setDeletingPostId(null);
    }
  };

  const handleUnrepost = (unrepostedPost: AggregatePost) => {
    // If this is a repost being removed, filter it out from the posts
    if (unrepostedPost.type === PostType.REPOST) {
      setPosts((prevPosts) =>
        prevPosts.filter(
          (p) =>
            !(
              p.type === PostType.REPOST &&
              p.repost?.id === unrepostedPost.repost?.id
            )
        )
      );
    }
  };

  const handleRepostCreated = (newRepost: AggregatePost) => {
    // When a new repost is created, refresh the feed to include it
    // This ensures the latest reposts appear in the feed
    handleRefresh();
  };

  const handlePostUpdate = (updatedPost: AggregatePost) => {
    setPosts((prevPosts) => {
      return prevPosts.map((post) => {
        const targetId =
          updatedPost.type === PostType.REPOST
            ? updatedPost.repost?.postId || updatedPost.id
            : updatedPost.id;

        const postTargetId =
          post.type === PostType.REPOST
            ? post.repost?.postId || post.id
            : post.id;

        // If this post references the same original content, sync interactions
        if (postTargetId === targetId) {
          return {
            ...post,
            liked: updatedPost.liked,
            reposted: updatedPost.reposted,
            favorited: updatedPost.favorited,
            likeCount: updatedPost.likeCount,
            favoriteCount: updatedPost.favoriteCount,
            repostCount: updatedPost.repostCount,
          };
        }

        return post;
      });
    });
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-8">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-heading font-bold text-primary">
            Your Feed
          </h2>
          <p className="text-muted-foreground">
            Latest posts from people you follow
          </p>
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={handleRefresh}
          disabled={isRefreshing}
          className="gap-2 bg-transparent"
        >
          <RefreshCw
            className={`h-4 w-4 ${isRefreshing ? "animate-spin" : ""}`}
          />
          Refresh
        </Button>
      </div>

      <CreatePost onPostCreated={handlePostCreated} auth={auth} />

      <div className="space-y-4">
        {posts.map((post, i) => (
          <PostCard
            key={getPostKey(post, i)}
            auth={auth}
            post={post}
            onEdit={handleEdit}
            onDelete={handleDelete}
            onPostUpdate={handlePostUpdate}
            onUnrepost={handleUnrepost}
            onRepostCreated={handleRepostCreated}
          />
        ))}
      </div>

      {pagination.hasMore && (
        <div className="flex justify-center pt-4">
          <Button
            onClick={handleLoadMore}
            disabled={isLoadingMore}
            variant="outline"
            className="gap-2 bg-transparent"
          >
            {isLoadingMore && <Loader2 className="h-4 w-4 animate-spin" />}
            Load More Posts
          </Button>
        </div>
      )}

      {posts.length === 0 && !isLoading && (
        <div className="text-center py-12">
          <div className="max-w-md mx-auto">
            <h3 className="text-lg font-heading font-semibold text-foreground mb-2">
              Your feed is empty
            </h3>
            <p className="text-muted-foreground mb-4">
              Follow some users to see their posts in your personalized feed, or
              create your first post!
            </p>
            <Button asChild>
              <a href="/search">Search People</a>
            </Button>
          </div>
        </div>
      )}

      <EditPostDialog
        auth={auth}
        post={editingPost}
        open={!!editingPost}
        onOpenChange={(open) => !open && setEditingPost(null)}
        onPostUpdated={handlePostUpdated}
      />

      <AlertDialog
        open={!!deletingPostId}
        onOpenChange={(open) => !open && setDeletingPostId(null)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Post</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete this post? This action cannot be
              undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={confirmDelete}
              className="bg-destructive text-destructive-foreground"
            >
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
