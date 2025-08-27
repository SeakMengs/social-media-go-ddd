"use client";

import { useState, useEffect } from "react";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import {
  Heart,
  Repeat2,
  Bookmark,
  MoreHorizontal,
  Edit,
  Trash2,
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { formatDistanceToNow } from "date-fns";
import { cn } from "@/lib/utils";
import { RepostDialog } from "./repost-dialog";
import { AggregatePost, Post } from "@/types/model";
import { AuthUser, AuthUserResult } from "@/auth";
import { toast } from "sonner";
import {
  favoritePost,
  likePost,
  unfavoritePost,
  unlikePost,
} from "@/api/action";

interface PostCardProps {
  auth: AuthUserResult;
  post: AggregatePost;
  onEdit?: (post: AggregatePost) => void;
  onDelete?: (postId: string) => void;
  onPostUpdate?: (post: AggregatePost) => void;
}

export function PostCard({
  post,
  onEdit,
  onDelete,
  onPostUpdate,
  auth,
}: PostCardProps) {
  const user = auth && auth.user;
  const [isLiked, setIsLiked] = useState(post.liked);
  const [isFavorited, setIsFavorited] = useState(post.favorited);
  const [isReposted, setIsReposted] = useState(false);
  const [showRepostDialog, setShowRepostDialog] = useState(false);
  const [isInteracting, setIsInteracting] = useState(false);
  const [likeCount, setLikeCount] = useState(post.likeCount);
  const [favoriteCount, setFavoriteCount] = useState(post.favoriteCount);
  const [repostCount, setRepostCount] = useState(post.repostCount);

  const isOwnPost = user?.id === post.userId;
  const timeAgo = post.createdAt
    ? formatDistanceToNow(new Date(post.createdAt), { addSuffix: true })
    : "Unknown";

  const handleLike = async () => {
    console.log(`handleLike called for post`, post);
    if (!user || isInteracting) return;

    const wasLiked = isLiked;
    // Optimistic update
    setIsLiked(!wasLiked);
    setLikeCount((prev) => prev + (wasLiked ? -1 : 1));

    setIsInteracting(true);
    try {
      const response = wasLiked
        ? await unlikePost(post.id)
        : await likePost(post.id);
      console.log(`like/unlike response: `, response);
      if (!response.success) {
        // Revert optimistic update on failure
        setIsLiked(wasLiked);
        setLikeCount((prev) => prev + (wasLiked ? 1 : -1));
        toast.error(response.error || "Failed to update like");
      } else {
        // Update parent component if needed
        onPostUpdate?.({
          ...post,
          likeCount: likeCount + (wasLiked ? -1 : 1),
        });
      }
    } catch (error) {
      // Revert optimistic update on error
      setIsLiked(wasLiked);
      setLikeCount((prev) => prev + (wasLiked ? 1 : -1));
      toast.error("An unexpected error occurred");
    } finally {
      setIsInteracting(false);
    }
  };

  const handleFavorite = async () => {
    if (!user || isInteracting) return;

    const wasFavorited = isFavorited;
    // Optimistic update
    setIsFavorited(!wasFavorited);
    setFavoriteCount((prev) => prev + (wasFavorited ? -1 : 1));

    setIsInteracting(true);
    try {
      const response = wasFavorited
        ? await unfavoritePost(post.id)
        : await favoritePost(post.id);

      if (!response.success) {
        // Revert optimistic update on failure
        setIsFavorited(wasFavorited);
        setFavoriteCount((prev) => prev + (wasFavorited ? 1 : -1));
        toast.error(response.error || "Failed to update favorite");
      } else {
        // Update parent component if needed
        onPostUpdate?.({
          ...post,
          favoriteCount: favoriteCount + (wasFavorited ? -1 : 1),
        });
      }
    } catch (error) {
      // Revert optimistic update on error
      setIsFavorited(wasFavorited);
      setFavoriteCount((prev) => prev + (wasFavorited ? 1 : -1));
      toast.error("An unexpected error occurred");
    } finally {
      setIsInteracting(false);
    }
  };

  const handleRepost = () => {
    if (!user || isInteracting) return;
    setShowRepostDialog(true);
  };

  const handleRepostCreated = (repostedPost: Post) => {
    setRepostCount((prev) => prev + 1);
    onPostUpdate?.({
      ...post,
      repostCount: repostCount + 1,
    });
  };

  const getInitials = (username: string) => {
    return username.slice(0, 2).toUpperCase();
  };

  return (
    <>
      <Card className="w-full shadow-sm hover:shadow-md transition-shadow duration-200">
        <CardHeader className="pb-3">
          <div className="flex items-start justify-between">
            <div className="flex items-center gap-3">
              <Avatar className="h-12 w-12 ring-2 ring-accent/20">
                <AvatarFallback className="bg-accent text-accent-foreground font-semibold">
                  {getInitials(post.user?.username || "U")}
                </AvatarFallback>
              </Avatar>
              <div className="flex flex-col">
                <div className="flex items-center gap-2">
                  <span className="font-heading font-bold text-foreground">
                    {post.user?.username}
                  </span>
                  {post.type === "repost" && (
                    <Badge
                      variant="secondary"
                      className="text-xs bg-accent/10 text-accent border-accent/20"
                    >
                      <Repeat2 className="h-3 w-3 mr-1" />
                      Reposted
                    </Badge>
                  )}
                </div>
                <span className="text-sm text-muted-foreground">{timeAgo}</span>
              </div>
            </div>
            {isOwnPost && (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-8 w-8 p-0 hover:bg-muted"
                  >
                    <MoreHorizontal className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem onClick={() => onEdit?.(post)}>
                    <Edit className="h-4 w-4 mr-2" />
                    Edit
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    onClick={() => onDelete?.(post.id)}
                    className="text-destructive"
                  >
                    <Trash2 className="h-4 w-4 mr-2" />
                    Delete
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            )}
          </div>
        </CardHeader>
        <CardContent className="pt-0">
          {post.type === "repost" && post.repost ? (
            <div className="space-y-4">
              {post.repost.comment && (
                <p className="text-foreground leading-relaxed text-pretty">
                  {post.repost.comment}
                </p>
              )}
              {post.user && (
                <Card className="border-l-4 border-l-accent bg-card/50 shadow-sm">
                  <CardContent className="p-4">
                    <div className="flex items-center gap-2 mb-3">
                      <Avatar className="h-6 w-6">
                        <AvatarFallback className="bg-accent text-accent-foreground text-xs">
                          {getInitials(post.user?.username || "U")}
                        </AvatarFallback>
                      </Avatar>
                      <span className="text-sm font-semibold">
                        {post.user?.username}
                      </span>
                    </div>
                    <p className="text-sm text-muted-foreground leading-relaxed text-pretty">
                      {post.content}
                    </p>
                  </CardContent>
                </Card>
              )}
            </div>
          ) : (
            <p className="text-foreground leading-relaxed text-pretty text-base">
              {post.content}
            </p>
          )}

          {user && (
            <div className="flex items-center justify-between mt-6 pt-4 border-t border-border">
              <div className="flex items-center gap-8">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={handleLike}
                  disabled={isInteracting}
                  className={cn(
                    "gap-2 text-muted-foreground hover:text-red-500 hover:bg-red-50 transition-all duration-200 rounded-full px-3",
                    isLiked && "text-red-500 bg-red-50"
                  )}
                >
                  <Heart
                    className={cn(
                      "h-4 w-4 transition-all",
                      isLiked && "fill-current scale-110"
                    )}
                  />
                  <span className="text-sm font-medium">{likeCount}</span>
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={handleRepost}
                  disabled={isInteracting}
                  className={cn(
                    "gap-2 text-muted-foreground hover:text-green-500 hover:bg-green-50 transition-all duration-200 rounded-full px-3",
                    isReposted && "text-green-500 bg-green-50"
                  )}
                >
                  <Repeat2
                    className={cn(
                      "h-4 w-4 transition-all",
                      isReposted && "scale-110"
                    )}
                  />
                  <span className="text-sm font-medium">{repostCount}</span>
                </Button>
              </div>
              <Button
                variant="ghost"
                size="sm"
                onClick={handleFavorite}
                disabled={isInteracting}
                className={cn(
                  "gap-2 text-muted-foreground hover:text-yellow-500 hover:bg-yellow-50 transition-all duration-200 rounded-full px-3",
                  isFavorited && "text-yellow-500 bg-yellow-50"
                )}
              >
                <Bookmark
                  className={cn(
                    "h-4 w-4 transition-all",
                    isFavorited && "fill-current scale-110"
                  )}
                />
                <span className="text-sm font-medium">{favoriteCount}</span>
              </Button>
            </div>
          )}
        </CardContent>
      </Card>

      {auth && (
        <RepostDialog
          auth={auth}
          post={post}
          open={showRepostDialog}
          onOpenChange={setShowRepostDialog}
          onRepostCreated={handleRepostCreated}
        />
      )}
    </>
  );
}
