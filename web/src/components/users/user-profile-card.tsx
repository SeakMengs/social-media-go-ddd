"use client";

import { useState } from "react";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import {
  CalendarDays,
  Users,
  FileText,
  UserPlus,
  UserMinus,
  Loader2,
} from "lucide-react";
import { formatDistanceToNow } from "date-fns";
import { AggregateUser } from "@/types/model";
import { AuthUserResult } from "@/auth";
import { followUser, unfollowUser } from "@/api/action";
import { toast } from "sonner";

interface UserProfileCardProps {
  auth: AuthUserResult;
  user: AggregateUser;
  postCount: number;
  repostCount: number;
  onFollowChange?: (userId: string, isFollowing: boolean) => void;
}

export function UserProfileCard({
  auth,
  user,
  postCount = 0,
  repostCount = 0,
  onFollowChange,
}: UserProfileCardProps) {
  const [isFollowing, setIsFollowing] = useState(user.followed);
  const [isLoading, setIsLoading] = useState(false);
  const currentUser = auth && auth.user;

  const isOwnProfile = currentUser?.id === user.id;
  const joinedDate = user.createdAt
    ? formatDistanceToNow(new Date(user.createdAt), { addSuffix: true })
    : "Unknown";

  const getInitials = (username: string) => {
    return username.slice(0, 2).toUpperCase();
  };

  const handleFollowToggle = async () => {
    if (!currentUser || isLoading) return;

    setIsLoading(true);
    const wasFollowing = isFollowing;
    // Optimistically update UI
    setIsFollowing(!wasFollowing);

    try {
      const response = wasFollowing
        ? await unfollowUser(user.id)
        : await followUser(user.id);

      if (!response.success) {
        // Revert UI on error
        setIsFollowing(wasFollowing);
        toast.error(response.error || "Failed to update follow status");
      } else {
        onFollowChange?.(user.id, !wasFollowing);
        toast.info(
          wasFollowing
            ? `You unfollowed ${user.username}`
            : `You are now following ${user.username}`
        );
      }
    } catch (error) {
      // Revert UI on error
      setIsFollowing(wasFollowing);
      toast.error("An unexpected error occurred");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card className="w-full">
      <CardHeader className="pb-4">
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-4">
            <Avatar className="h-16 w-16">
              <AvatarFallback className="bg-accent text-accent-foreground text-lg">
                {getInitials(user.username)}
              </AvatarFallback>
            </Avatar>
            <div className="space-y-1">
              <h2 className="text-xl font-heading font-bold text-foreground">
                {user.username}
              </h2>
              <p className="text-sm text-muted-foreground">{user.email}</p>
              <div className="flex items-center gap-2 text-xs text-muted-foreground">
                <CalendarDays className="h-3 w-3" />
                <span>Joined {joinedDate}</span>
              </div>
            </div>
          </div>
          {!isOwnProfile && currentUser && (
            <Button
              onClick={handleFollowToggle}
              disabled={isLoading}
              variant={isFollowing ? "outline" : "default"}
              className="gap-2"
            >
              {isLoading ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : isFollowing ? (
                <>
                  <UserMinus className="h-4 w-4" />
                  Unfollow
                </>
              ) : (
                <>
                  <UserPlus className="h-4 w-4" />
                  Follow
                </>
              )}
            </Button>
          )}
          {isOwnProfile && (
            <Badge variant="secondary" className="gap-1">
              <Users className="h-3 w-3" />
              Your Profile
            </Badge>
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex items-center gap-6 text-sm">
          <div className="flex items-center gap-1">
            <Users className="h-4 w-4 text-muted-foreground" />
            <span className="font-semibold text-foreground">
              {user.followerCount.toLocaleString()}
            </span>
            <span className="text-muted-foreground">followers</span>
          </div>
          <div className="flex items-center gap-1">
            <Users className="h-4 w-4 text-muted-foreground" />
            <span className="font-semibold text-foreground">
              {user.followingCount.toLocaleString()}
            </span>
            <span className="text-muted-foreground">following</span>
          </div>
          <div className="flex items-center gap-1">
            <FileText className="h-4 w-4 text-muted-foreground" />
            <span className="font-semibold text-foreground">
              {postCount.toLocaleString()}
            </span>
            <span className="text-muted-foreground">posts</span>
          </div>
          <div className="flex items-center gap-1">
            <FileText className="h-4 w-4 text-muted-foreground" />
            <span className="font-semibold text-foreground">
              {repostCount.toLocaleString()}
            </span>
            <span className="text-muted-foreground">reposts</span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
