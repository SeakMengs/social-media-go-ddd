"use client";

import type React from "react";

import { useState, useEffect } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Card, CardContent } from "@/components/ui/card";
import { Loader2, Repeat2 } from "lucide-react";
import { AggregatePost, PostType } from "@/types/model";
import { AuthUser } from "@/auth";
import { repost, unrepost } from "@/api/action";
import { toast } from "sonner";
import { getPostId } from "@/utils/post";

interface RepostDialogProps {
  auth: AuthUser;
  post: AggregatePost | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onRepostCreated?: (post: AggregatePost) => void;
  onRepostDeleted?: (post: AggregatePost) => void;
}

export function RepostDialog({
  auth,
  post,
  open,
  onOpenChange,
  onRepostCreated,
  onRepostDeleted,
}: RepostDialogProps) {
  const [comment, setComment] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const { user } = auth;

  // Set default comment when dialog opens for updating a repost
  useEffect(() => {
    if (open && post?.reposted) {
      // If this is a repost being viewed and user has reposted it, get the comment
      if (post.type === PostType.REPOST && post.repost?.comment) {
        setComment(post.repost.comment);
      } else {
        // For regular posts that user has reposted, we don't have the comment here
        // The backend should provide this when opening the dialog
        setComment("");
      }
    } else if (open) {
      setComment("");
    }
  }, [post, open]);

  const getInitials = (username: string) => {
    return username.slice(0, 2).toUpperCase();
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!post || !user) return;

    setIsLoading(true);
    try {
      // For reposts, target the original post, not the repost
      const targetPostId = getPostId(post);
      const response = await repost(targetPostId, {
        comment: comment.trim(),
      });
      if (response.success) {
        onOpenChange(false);
        setComment("");
        onRepostCreated?.(post);
        toast.success(post.reposted ? "Repost updated successfully!" : "Post reposted successfully!");
      } else {
        toast.error(response.error || "Failed to repost");
      }
    } catch (error) {
      toast.error("An unexpected error occurred");
    } finally {
      setIsLoading(false);
    }
  };

  const handleUnrepost = async () => {
    if (!post || !user) return;

    setIsLoading(true);
    try {
      const targetPostId = getPostId(post);
      const response = await unrepost(targetPostId);
      
      if (response.success) {
        onOpenChange(false);
        setComment("");
        onRepostDeleted?.(post);
        toast.success("Repost removed successfully!");
      } else {
        toast.error(response.error || "Failed to remove repost");
      }
    } catch (error) {
      toast.error("An unexpected error occurred");
    } finally {
      setIsLoading(false);
    }
  };

  if (!post || !user) return null;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[525px]">
        <DialogHeader>
          <DialogTitle className="font-heading flex items-center gap-2">
            <Repeat2 className="h-5 w-5" />
            {post.reposted ? "Update Repost" : "Repost"}
          </DialogTitle>
          <DialogDescription>
            {post.reposted 
              ? "Update your repost comment or remove your repost."
              : "Add your thoughts or repost as-is"
            }
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          <div className="flex gap-3">
            <Avatar className="h-10 w-10 flex-shrink-0">
              <AvatarFallback className="bg-accent text-accent-foreground">
                {getInitials(user.username)}
              </AvatarFallback>
            </Avatar>
            <div className="flex-1">
              <Textarea
                placeholder="Add a comment (optional)"
                value={comment}
                onChange={(e) => setComment(e.target.value)}
                className="min-h-[80px] resize-none"
                disabled={isLoading}
              />
            </div>
          </div>

          <Card className="border-l-4 border-l-accent bg-muted/30">
            <CardContent className="p-4">
              <div className="flex items-center gap-2 mb-2">
                <Avatar className="h-6 w-6">
                  <AvatarFallback className="bg-accent text-accent-foreground text-xs">
                    {getInitials(post.user?.username || "U")}
                  </AvatarFallback>
                </Avatar>
                <span className="text-sm font-medium">
                  {post.user?.username}
                </span>
              </div>
              <p className="text-sm text-muted-foreground leading-relaxed">
                {post.content}
              </p>
            </CardContent>
          </Card>
        </div>

        <DialogFooter className="gap-2">
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={isLoading}
          >
            Cancel
          </Button>
          {post.reposted && (
            <Button
              type="button"
              variant="destructive"
              onClick={handleUnrepost}
              disabled={isLoading}
            >
              {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Remove Repost
            </Button>
          )}
          <Button onClick={handleSubmit} disabled={isLoading}>
            {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {post.reposted ? "Update Repost" : "Repost"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
