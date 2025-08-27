"use client";

import { Suspense, useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Search, Users, UserPlus, UserCheck } from "lucide-react";
import Link from "next/link";
import useDebounce from "@/hooks/use-debounce";
import { AggregateUser } from "@/types/model";
import { AuthUser, AuthUserResult } from "@/auth";
import { followUser, searchUsersByName, unfollowUser } from "@/api/action";
import { toast } from "sonner";

type SearchContentProps = {
  auth: AuthUserResult;
};

export function SearchContent({ auth }: SearchContentProps) {
  const searchParams = useSearchParams();
  const initialQuery = searchParams.get("q") || "";
  const [searchQuery, setSearchQuery] = useState(initialQuery);
  const [searchResults, setSearchResults] = useState<AggregateUser[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [followingUsers, setFollowingUsers] = useState<Set<string>>(new Set());
  const debouncedQuery = useDebounce(searchQuery, 300);

  useEffect(() => {
    if (debouncedQuery.trim()) {
      handleSearch(debouncedQuery);
    } else {
      setSearchResults([]);
    }
  }, [debouncedQuery]);

  const handleSearch = async (query: string) => {
    if (!query.trim()) return;

    setIsLoading(true);
    try {
      const results = await searchUsersByName(query);
      console.log("search by users name", results);
      if (results.success) {
        setSearchResults(results.data.users);
      } else {
        toast.error(results.error || "Failed to fetch search results");
      }
    } catch (error) {
      if (error instanceof Error) {
        toast.error(error.message || "An unexpected error occurred");
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleFollow = async (targetUserId: string) => {
    if (!auth || !auth.user) return;

    // Find the user in searchResults to get the current followed status
    const searchUser = searchResults.find((u) => u.id === targetUserId);
    if (!searchUser) return;

    const wasFollowing = searchUser.followed;
    // Optimistically update UI
    setSearchResults((prev) =>
      prev.map((u) =>
        u.id === targetUserId
          ? {
              ...u,
              followed: !wasFollowing,
              followerCount: wasFollowing
                ? u.followerCount - 1
                : u.followerCount + 1,
            }
          : u
      )
    );

    try {
      const response = wasFollowing
        ? await unfollowUser(targetUserId)
        : await followUser(targetUserId);

      if (!response.success) {
        // Revert UI on error
        setSearchResults((prev) =>
          prev.map((u) =>
            u.id === targetUserId
              ? {
                  ...u,
                  followed: wasFollowing,
                  followerCount: wasFollowing
                    ? u.followerCount + 1
                    : u.followerCount - 1,
                }
              : u
          )
        );
        toast.error(response.error || "Failed to update follow status");
      }
    } catch (error) {
      // Revert UI on error
      setSearchResults((prev) =>
        prev.map((u) =>
          u.id === targetUserId
            ? {
                ...u,
                followed: wasFollowing,
                followerCount: wasFollowing
                  ? u.followerCount + 1
                  : u.followerCount - 1,
              }
            : u
        )
      );
      toast.error("An unexpected error occurred");
    }
  };

  const getInitials = (username: string) => {
    return username.slice(0, 2).toUpperCase();
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="space-y-6">
        <div className="text-center space-y-2">
          <h1 className="font-heading font-black text-3xl text-primary">
            Search Users
          </h1>
          <p className="text-muted-foreground">
            Find and connect with other users
          </p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Search className="h-5 w-5" />
              Search
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                type="search"
                placeholder="Search for users..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10 bg-input border-border focus:ring-accent focus:border-accent"
              />
            </div>
          </CardContent>
        </Card>

        {isLoading && (
          <div className="text-center py-8">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-accent mx-auto"></div>
            <p className="text-muted-foreground mt-2">Searching...</p>
          </div>
        )}

        {Array.isArray(searchResults) && searchResults.length > 0 && (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Users className="h-5 w-5" />
                Search Results ({searchResults.length})
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {searchResults.map((searchUser) => (
                  <div
                    key={searchUser.id}
                    className="flex items-center justify-between p-4 border border-border rounded-lg hover:bg-muted/50 transition-colors"
                  >
                    <Link
                      href={`/profile/${searchUser.id}`}
                      className="flex items-center gap-4 flex-1"
                    >
                      <Avatar className="h-12 w-12 ring-2 ring-accent/20">
                        <AvatarFallback className="bg-accent text-accent-foreground font-semibold">
                          {getInitials(searchUser.username)}
                        </AvatarFallback>
                      </Avatar>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 mb-1">
                          <h3 className="font-heading font-bold text-foreground truncate">
                            {searchUser.username}
                          </h3>
                        </div>
                        {/* {searchUser.bio && (
                      <p className="text-sm text-muted-foreground line-clamp-2 text-pretty">{searchUser.bio}</p>
                    )} */}
                        <div className="flex items-center gap-4 mt-2 text-xs text-muted-foreground">
                          <span>{searchUser.followerCount} followers</span>
                          {/* <span>{searchUser.postCount} posts</span> */}
                        </div>
                      </div>
                    </Link>
                    {auth &&
                      (searchUser.id === auth.user.id ? (
                        <Badge variant="outline" className="ml-4">
                          You
                        </Badge>
                      ) : (
                        <Button
                          onClick={() => handleFollow(searchUser.id)}
                          variant={
                            searchUser.followed ? "outline" : "default"
                          }
                          size="sm"
                          className="ml-4"
                        >
                          {searchUser.followed ? (
                            <>
                              <UserCheck className="h-4 w-4 mr-2" />
                              Following
                            </>
                          ) : (
                            <>
                              <UserPlus className="h-4 w-4 mr-2" />
                              Follow
                            </>
                          )}
                        </Button>
                      ))}
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        )}

        {searchQuery &&
          !isLoading &&
          Array.isArray(searchResults) &&
          searchResults.length === 0 && (
            <Card>
              <CardContent className="text-center py-8">
                <Users className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
                <h3 className="font-heading font-bold text-lg mb-2">
                  No users found
                </h3>
                <p className="text-muted-foreground">
                  Try searching with different keywords
                </p>
              </CardContent>
            </Card>
          )}
      </div>
    </div>
  );
}
