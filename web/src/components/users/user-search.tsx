"use client"

import { useState, useEffect } from "react"
import { Input } from "@/components/ui/input"
import { Card, CardContent } from "@/components/ui/card"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { Search, Loader2, UserPlus } from "lucide-react"
import { searchUsers, followUser } from "@/lib/users"
import type { UserProfile } from "@/lib/users"
import { useAuth } from "@/hooks/use-auth"
import { useToast } from "@/hooks/use-toast"
import { useDebounce } from "@/hooks/use-debounce"
import Link from "next/link"

export function UserSearch() {
  const [query, setQuery] = useState("")
  const [users, setUsers] = useState<UserProfile[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [followingUsers, setFollowingUsers] = useState<Set<string>>(new Set())
  const debouncedQuery = useDebounce(query, 300)
  const { user: currentUser } = useAuth()
  const { toast } = useToast()

  useEffect(() => {
    const performSearch = async () => {
      if (!debouncedQuery.trim()) {
        setUsers([])
        return
      }

      setIsLoading(true)
      try {
        const response = await searchUsers(debouncedQuery)
        if (response.success && response.data?.users) {
          setUsers(response.data.users.filter((user) => user.id !== currentUser?.id))
        }
      } catch (error) {
        toast({
          title: "Error",
          description: "Failed to search users",
          variant: "destructive",
        })
      } finally {
        setIsLoading(false)
      }
    }

    performSearch()
  }, [debouncedQuery, currentUser?.id, toast])

  const getInitials = (username: string) => {
    return username.slice(0, 2).toUpperCase()
  }

  const handleFollow = async (userId: string) => {
    if (!currentUser) return

    setFollowingUsers((prev) => new Set(prev).add(userId))
    try {
      const response = await followUser(userId, currentUser.id)
      if (response.success) {
        setUsers((prev) => prev.map((user) => (user.id === userId ? { ...user, isFollowing: true } : user)))
        toast({
          title: "Success",
          description: "User followed successfully",
        })
      }
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to follow user",
        variant: "destructive",
      })
    } finally {
      setFollowingUsers((prev) => {
        const newSet = new Set(prev)
        newSet.delete(userId)
        return newSet
      })
    }
  }

  return (
    <div className="space-y-4">
      <div className="relative">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="Search users..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="pl-10"
        />
      </div>

      {isLoading && (
        <div className="flex items-center justify-center py-8">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </div>
      )}

      {users.length > 0 && (
        <div className="space-y-3">
          {users.map((user) => (
            <Card key={user.id} className="hover:bg-muted/50 transition-colors">
              <CardContent className="p-4">
                <div className="flex items-center justify-between">
                  <Link href={`/profile/${user.id}`} className="flex items-center gap-3 flex-1 hover:opacity-80">
                    <Avatar className="h-10 w-10">
                      <AvatarFallback className="bg-accent text-accent-foreground">
                        {getInitials(user.username)}
                      </AvatarFallback>
                    </Avatar>
                    <div className="flex-1 min-w-0">
                      <h3 className="font-semibold text-foreground truncate">{user.username}</h3>
                      <p className="text-sm text-muted-foreground truncate">{user.bio || user.email}</p>
                      <div className="flex items-center gap-4 text-xs text-muted-foreground mt-1">
                        <span>{user.followerCount} followers</span>
                        <span>{user.postCount} posts</span>
                      </div>
                    </div>
                  </Link>
                  {!user.isFollowing && (
                    <Button
                      size="sm"
                      onClick={() => handleFollow(user.id)}
                      disabled={followingUsers.has(user.id)}
                      className="gap-2 ml-3"
                    >
                      {followingUsers.has(user.id) ? (
                        <Loader2 className="h-3 w-3 animate-spin" />
                      ) : (
                        <UserPlus className="h-3 w-3" />
                      )}
                      Follow
                    </Button>
                  )}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {query && !isLoading && users.length === 0 && (
        <div className="text-center py-8 text-muted-foreground">
          <Search className="h-12 w-12 mx-auto mb-4 opacity-50" />
          <p>No users found for "{query}"</p>
        </div>
      )}
    </div>
  )
}
