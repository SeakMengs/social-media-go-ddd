"use client";

import type React from "react";

import { useState } from "react";
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
import { AggregatePost } from "@/types/model";
import { AuthUser } from "@/auth";
import { repost } from "@/api/action";
import { toast } from "sonner";

interface RepostDialogProps {
  auth: AuthUser;
  post: AggregatePost | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onRepostCreated?: (post: AggregatePost) => void;
}

export function RepostDialog({
  auth,
  post,
  open,
  onOpenChange,
  onRepostCreated,
}: RepostDialogProps) {
  const [comment, setComment] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const { user } = auth;

  const getInitials = (username: string) => {
    return username.slice(0, 2).toUpperCase();
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!post || !user) return;

    setIsLoading(true);
    try {
      const response = await repost(post.id, {
        comment: comment.trim(),
      });
      if (response.success) {
        onOpenChange(false);
        setComment("");
        onRepostCreated?.(post);
        toast.success("Post reposted successfully!");
      } else {
        toast.error(response.error || "Failed to repost");
      }
    } catch (error) {
      toast.error("An unexpected error occurred");
    } finally {
      setIsLoading(false);
    }
  };

  const handleQuickRepost = async () => {
    if (!post || !user) return;

    setIsLoading(true);
    try {
      const response = await repost(post.id, {
        comment: "",
      });
      if (response.success) {
        onOpenChange(false);
        onRepostCreated?.(post);
        toast.success("Post reposted successfully!");
      } else {
        toast.error(response.error || "Failed to repost");
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
            Repost
          </DialogTitle>
          <DialogDescription>
            Add your thoughts or repost as-is
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
          <Button
            type="button"
            variant="outline"
            onClick={handleQuickRepost}
            disabled={isLoading}
          >
            {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Repost
          </Button>
          <Button onClick={handleSubmit} disabled={isLoading}>
            {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Repost with comment
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
