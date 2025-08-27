"use client"

import { useState, useEffect } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { UserPlus, Loader2 } from "lucide-react"
import { getSuggestedUsers, followUser } from "@/lib/users"
import type { UserProfile } from "@/lib/users"
import { useAuth } from "@/hooks/use-auth"
import { useToast } from "@/hooks/use-toast"
import Link from "next/link"

export function SuggestedUsers() {
  const [users, setUsers] = useState<UserProfile[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [followingUsers, setFollowingUsers] = useState<Set<string>>(new Set())
  const { user: currentUser } = useAuth()
  const { toast } = useToast()

  useEffect(() => {
    const fetchSuggestedUsers = async () => {
      if (!currentUser) return

      try {
        const response = await getSuggestedUsers(currentUser.id)
        if (response.success && response.data?.users) {
          setUsers(response.data.users)
        }
      } catch (error) {
        toast({
          title: "Error",
          description: "Failed to load suggested users",
          variant: "destructive",
        })
      } finally {
        setIsLoading(false)
      }
    }

    fetchSuggestedUsers()
  }, [currentUser, toast])

  const getInitials = (username: string) => {
    return username.slice(0, 2).toUpperCase()
  }

  const handleFollow = async (userId: string) => {
    if (!currentUser) return

    setFollowingUsers((prev) => new Set(prev).add(userId))
    try {
      const response = await followUser(userId, currentUser.id)
      if (response.success) {
        setUsers((prev) => prev.filter((user) => user.id !== userId))
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

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="text-lg font-heading">Suggested for you</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-center py-4">
            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
          </div>
        </CardContent>
      </Card>
    )
  }

  if (users.length === 0) {
    return null
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-lg font-heading">Suggested for you</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {users.map((user) => (
          <div key={user.id} className="flex items-center justify-between">
            <Link href={`/profile/${user.id}`} className="flex items-center gap-3 flex-1 hover:opacity-80">
              <Avatar className="h-10 w-10">
                <AvatarFallback className="bg-accent text-accent-foreground">
                  {getInitials(user.username)}
                </AvatarFallback>
              </Avatar>
              <div className="flex-1 min-w-0">
                <h4 className="font-semibold text-foreground truncate">{user.username}</h4>
                <p className="text-sm text-muted-foreground truncate">
                  {user.followerCount} followers • {user.postCount} posts
                </p>
              </div>
            </Link>
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
          </div>
        ))}
      </CardContent>
    </Card>
  )
}
