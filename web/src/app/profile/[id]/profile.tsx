"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { UserProfileCard } from "@/components/users/user-profile-card";
import { PostCard } from "@/components/posts/post-card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Loader2 } from "lucide-react";
import { AggregatePost, AggregateUser } from "@/types/model";
import { AuthUserResult } from "@/auth";
import { getUserById, getUserPosts, getUserReposts } from "@/api/action";
import { toast } from "sonner";
import { getPostKey } from "@/utils/key";

type ProfileProps = {
  auth: AuthUserResult;
};

export default function Profile({ auth }: ProfileProps) {
  const params = useParams();
  const [user, setUser] = useState<AggregateUser | null>(null);
  const [posts, setPosts] = useState<AggregatePost[]>([]);
  const [reposts, setReposts] = useState<AggregatePost[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [postsLoading, setPostsLoading] = useState(true);
  const [repostsLoading, setRepostsLoading] = useState(true);

  const userId = params.id as string;

  useEffect(() => {
    const fetchUserProfile = async () => {
      try {
        const response = await getUserById(userId);
        console.log(`user ${userId} response: `, response);
        if (response.success) {
          setUser(response.data.user);
        } else {
          toast.error(response.error || "Failed to get user information");
        }
      } catch (error) {
        if (error instanceof Error) {
          toast.error(error.message || "An unexpected error occurred");
        }
      } finally {
        setIsLoading(false);
      }
    };

    fetchUserProfile();
  }, [userId]);

  useEffect(() => {
    const fetchUserPosts = async () => {
      try {
        const response = await getUserPosts(userId);
        console.log(`user ${userId} posts`, response);
        if (response.success) {
          setPosts(response.data.posts);
        } else {
          toast.error(response.error || "Failed to get user's posts");
        }
      } catch (error) {
        if (error instanceof Error) {
          toast.error(error.message || "An unexpected error occurred");
        }
      } finally {
        setPostsLoading(false);
      }
    };
    const fetchUserReposts = async () => {
      try {
        const response = await getUserReposts(userId);
        console.log(`user ${userId} reposts`, response);
        if (response.success) {
          setReposts(response.data.reposts);
        } else {
          toast.error(response.error || "Failed to get user's posts");
        }
      } catch (error) {
        if (error instanceof Error) {
          toast.error(error.message || "An unexpected error occurred");
        }
      } finally {
        setRepostsLoading(false);
      }
    };

    if (userId) {
      Promise.all([fetchUserPosts(), fetchUserReposts()]);
    }
  }, [userId]);

  const handleFollowChange = (userId: string, isFollowing: boolean) => {
    if (user) {
      setUser({
        ...user,
        followed: isFollowing,
        followerCount: user.followerCount + (isFollowing ? 1 : -1),
      });
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (!user) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-heading font-bold text-foreground mb-2">
            User not found
          </h1>
          <p className="text-muted-foreground">
            The user you're looking for doesn't exist.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <div className="container max-w-4xl mx-auto py-8 px-4">
        <div className="space-y-6">
          <UserProfileCard
            user={user}
            onFollowChange={handleFollowChange}
            postCount={posts.length}
            repostCount={reposts.length}
            auth={auth}
          />

          <Tabs defaultValue="posts" className="w-full">
            <TabsList className="grid w-full grid-cols-2">
              <TabsTrigger value="posts">Posts ({posts.length})</TabsTrigger>
              <TabsTrigger value="reposts">Reposts</TabsTrigger>
            </TabsList>
            <TabsContent value="posts" className="space-y-4 mt-6">
              {postsLoading ? (
                <div className="flex items-center justify-center py-8">
                  <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                </div>
              ) : posts.length > 0 ? (
                posts.map((post) => (
                  <PostCard key={getPostKey(post)} post={post} auth={auth} />
                ))
              ) : (
                <div className="text-center py-8 text-muted-foreground">
                  <p>No posts yet</p>
                </div>
              )}
            </TabsContent>
            <TabsContent value="reposts" className="space-y-4 mt-6">
              {repostsLoading ? (
                <div className="flex items-center justify-center py-8">
                  <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                </div>
              ) : reposts.length > 0 ? (
                reposts.map((repost) => (
                  <PostCard key={getPostKey(repost)} post={repost} auth={auth} />
                ))
              ) : (
                <div className="text-center py-8 text-muted-foreground">
                  <p>No reposts yet</p>
                </div>
              )}
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </div>
  );
}
