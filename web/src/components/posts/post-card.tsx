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
import { AggregatePost, Post, PostType } from "@/types/model";
import { AuthUser, AuthUserResult } from "@/auth";
import { toast } from "sonner";
import {
  favoritePost,
  likePost,
  unfavoritePost,
  unlikePost,
  unrepost,
} from "@/api/action";
import { getPostId } from "@/utils/post";

interface PostCardProps {
  auth: AuthUserResult;
  post: AggregatePost;
  onEdit?: (post: AggregatePost) => void;
  onDelete?: (postId: string) => void;
  onPostUpdate?: (post: AggregatePost) => void;
  onUnrepost?: (post: AggregatePost) => void;
  onRepostCreated?: (post: AggregatePost) => void;
}

export function PostCard({
  post,
  onEdit,
  onDelete,
  onPostUpdate,
  onUnrepost,
  onRepostCreated,
  auth,
}: PostCardProps) {
  const user = auth && auth.user;
  const [isLiked, setIsLiked] = useState(post.liked);
  const [isFavorited, setIsFavorited] = useState(post.favorited);
  const [isReposted, setIsReposted] = useState(post.reposted);
  const [showRepostDialog, setShowRepostDialog] = useState(false);
  const [isInteracting, setIsInteracting] = useState(false);
  const [likeCount, setLikeCount] = useState(post.likeCount);
  const [favoriteCount, setFavoriteCount] = useState(post.favoriteCount);
  const [repostCount, setRepostCount] = useState(post.repostCount);

  // Sync local state with prop changes (for updates from parent components)
  useEffect(() => {
    setIsLiked(post.liked);
    setIsFavorited(post.favorited);
    setIsReposted(post.reposted);
    setLikeCount(post.likeCount);
    setFavoriteCount(post.favoriteCount);
    setRepostCount(post.repostCount);
  }, [post.liked, post.favorited, post.reposted, post.likeCount, post.favoriteCount, post.repostCount]);

  const isOwnPost = post.type === PostType.REPOST 
    ? user?.id === post.repostUser?.id  // For reposts, check if current user is the reposter
    : user?.id === post.userId;         // For regular posts, check if current user is the author
  const timeAgo = post.createdAt
    ? formatDistanceToNow(new Date(post.createdAt), { addSuffix: true })
    : "Unknown";

  const handleLike = async () => {
    console.log(`handleLike called for post`, post);
    if (!user || isInteracting) return;

    const wasLiked = isLiked;
    const newLikeCount = likeCount + (wasLiked ? -1 : 1);
    
    // Optimistic update
    setIsLiked(!wasLiked);
    setLikeCount(newLikeCount);

    setIsInteracting(true);
    try {
      // For reposts, interact with the original post, not the repost
      const targetPostId = getPostId(post);
      const response = wasLiked
        ? await unlikePost(targetPostId)
        : await likePost(targetPostId);
      console.log(`like/unlike response: `, response);
      if (!response.success) {
        // Revert optimistic update on failure
        setIsLiked(wasLiked);
        setLikeCount(likeCount);
        toast.error(response.error || "Failed to update like");
      } else {
        // Update parent component with the new state
        onPostUpdate?.({
          ...post,
          liked: !wasLiked,
          likeCount: newLikeCount,
        });
      }
    } catch (error) {
      // Revert optimistic update on error
      setIsLiked(wasLiked);
      setLikeCount(likeCount);
      toast.error("An unexpected error occurred");
    } finally {
      setIsInteracting(false);
    }
  };

    const handleFavorite = async () => {
    if (!user || isInteracting) return;

    const wasFavorited = isFavorited;
    const newFavoriteCount = favoriteCount + (wasFavorited ? -1 : 1);
    
    // Optimistic update
    setIsFavorited(!wasFavorited);
    setFavoriteCount(newFavoriteCount);

    setIsInteracting(true);
    try {
      // For reposts, interact with the original post, not the repost
      const targetPostId = getPostId(post);
      const response = wasFavorited
        ? await unfavoritePost(targetPostId)
        : await favoritePost(targetPostId);

      if (!response.success) {
        // Revert optimistic update on failure
        setIsFavorited(wasFavorited);
        setFavoriteCount(favoriteCount);
        toast.error(response.error || "Failed to update favorite");
      } else {
        // Update parent component with the new state
        onPostUpdate?.({
          ...post,
          favorited: !wasFavorited,
          favoriteCount: newFavoriteCount,
        });
      }
    } catch (error) {
      // Revert optimistic update on error
      setIsFavorited(wasFavorited);
      setFavoriteCount(favoriteCount);
      toast.error("An unexpected error occurred");
    } finally {
      setIsInteracting(false);
    }
  };

  const handleRepost = async () => {
    if (!user || isInteracting) return;
    
    setShowRepostDialog(true);
  };

  const handleRepostCreated = (repostedPost: AggregatePost) => {
    const newRepostCount = repostCount + 1;
    setRepostCount(newRepostCount);
    setIsReposted(true);
    
    const updatedPost = {
      ...post,
      repostCount: newRepostCount,
      reposted: true,
    };
    
    onPostUpdate?.(updatedPost);
    onRepostCreated?.(updatedPost);
  };

  const handleRepostDeleted = (deletedPost: AggregatePost) => {
    const newRepostCount = repostCount - 1;
    setRepostCount(newRepostCount);
    setIsReposted(false);
    
    const updatedPost = {
      ...post,
      repostCount: newRepostCount,
      reposted: false,
    };
    
    onPostUpdate?.(updatedPost);
    onUnrepost?.(updatedPost);
  };

  const getInitials = (username: string) => {
    return username.slice(0, 2).toUpperCase();
  };

  return (
    <>
      {post.type === PostType.REPOST && post.repost ? (
        // Repost Layout - Different structure to make it clear this is a repost
        <Card className="w-full shadow-sm hover:shadow-md transition-shadow duration-200">
          <CardContent className="p-4 space-y-3">
            {/* Repost Header */}
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Repeat2 className="h-4 w-4" />
              <Avatar className="h-6 w-6">
                <AvatarFallback className="bg-accent text-accent-foreground text-xs">
                  {getInitials(post.repostUser?.username || "U")}
                </AvatarFallback>
              </Avatar>
              <span className="font-medium">{post.repostUser?.username}</span>
              <span>reposted</span>
              <span>{timeAgo}</span>
            </div>
            
            {/* Repost Comment */}
            {post.repost.comment && (
              <div className="pl-6">
                <p className="text-foreground leading-relaxed text-pretty">
                  {post.repost.comment}
                </p>
              </div>
            )}
            
            {/* Original Post Card */}
            <Card className="shadow-sm border-l-4 border-l-accent bg-card/50 ml-6">
              <CardHeader className="pb-3">
                <div className="flex items-start justify-between">
                  <div className="flex items-center gap-3">
                    <Avatar className="h-10 w-10 ring-2 ring-accent/20">
                      <AvatarFallback className="bg-accent text-accent-foreground font-semibold">
                        {getInitials(post.user?.username || "U")}
                      </AvatarFallback>
                    </Avatar>
                    <div className="flex flex-col">
                      <span className="font-heading font-bold text-foreground">
                        {post.user?.username}
                      </span>
                      <span className="text-sm text-muted-foreground">Original post</span>
                    </div>
                  </div>
                </div>
              </CardHeader>
              <CardContent className="pt-0">
                <p className="text-foreground leading-relaxed text-pretty text-base mb-4">
                  {post.content}
                </p>
                
                {user && (
                  <div className="flex items-center justify-between pt-4 border-t border-border">
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
          </CardContent>
        </Card>
      ) : (
        // Regular Post Layout
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
                  <span className="font-heading font-bold text-foreground">
                    {post.user?.username}
                  </span>
                  <span className="text-sm text-muted-foreground">{timeAgo}</span>
                </div>
              </div>
              {(isOwnPost || (post.type === PostType.REPOST && user?.id === post.repostUser?.id)) && (
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
                    {post.type === PostType.REPOST && user?.id === post.repostUser?.id ? (
                      // Menu for user's own repost
                      <>
                        <DropdownMenuItem onClick={() => setShowRepostDialog(true)}>
                          <Edit className="h-4 w-4 mr-2" />
                          Edit Repost
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          onClick={async () => {
                            setIsInteracting(true);
                            try {
                              const targetPostId = getPostId(post);
                              const response = await unrepost(targetPostId);
                              if (response.success) {
                                const updatedPost = {
                                  ...post,
                                  reposted: false,
                                  repostCount: repostCount - 1,
                                };
                                handleRepostDeleted(updatedPost);
                                toast.success("Repost removed successfully!");
                              } else {
                                toast.error(response.error || "Failed to remove repost");
                              }
                            } catch (error) {
                              toast.error("An unexpected error occurred");
                            } finally {
                              setIsInteracting(false);
                            }
                          }}
                          className="text-destructive"
                        >
                          <Trash2 className="h-4 w-4 mr-2" />
                          Remove Repost
                        </DropdownMenuItem>
                      </>
                    ) : (
                      // Menu for user's own original post
                      <>
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
                      </>
                    )}
                  </DropdownMenuContent>
                </DropdownMenu>
              )}
            </div>
          </CardHeader>
          <CardContent className="pt-0">
            <p className="text-foreground leading-relaxed text-pretty text-base">
              {post.content}
            </p>

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
      )}

      {auth && (
        <RepostDialog
          auth={auth}
          post={post}
          open={showRepostDialog}
          onOpenChange={setShowRepostDialog}
          onRepostCreated={handleRepostCreated}
          onRepostDeleted={handleRepostDeleted}
        />
      )}
    </>
  );
}
