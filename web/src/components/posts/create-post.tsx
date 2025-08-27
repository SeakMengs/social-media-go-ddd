"use client"

import type React from "react"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Textarea } from "@/components/ui/textarea"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Loader2, Send } from "lucide-react"
import { AuthUser } from "@/auth"
import { createPost } from "@/api/action"
import { toast } from "sonner"

interface CreatePostProps {
  auth: AuthUser
  onPostCreated?: () => void
}

export function CreatePost({ auth, onPostCreated }: CreatePostProps) {
  const [content, setContent] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const { user } = auth

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!content.trim() || !user) return

    setIsLoading(true)
    try {
      const response = await createPost({
        content: content.trim(),
      })
      if (response.success) {
        setContent("")
        onPostCreated?.()
        toast.success("Post created successfully")
      } else {
        toast.error(response.error || "Failed to create post")
      }
    } catch (error) {
      toast.error("An unexpected error occurred")
    } finally {
      setIsLoading(false)
    }
  }

  const getInitials = (username: string) => {
    return username.slice(0, 2).toUpperCase()
  }

  if (!user) return null

  return (
    <Card className="w-full">
      <CardHeader className="pb-3">
        <CardTitle className="text-lg font-heading">Share your thoughts</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="flex gap-3">
            <Avatar className="h-10 w-10 flex-shrink-0">
              <AvatarFallback className="bg-accent text-accent-foreground">{getInitials(user.username)}</AvatarFallback>
            </Avatar>
            <div className="flex-1">
              <Textarea
                placeholder="What's on your mind?"
                value={content}
                onChange={(e) => setContent(e.target.value)}
                className="min-h-[100px] resize-none border-0 p-0 text-base placeholder:text-muted-foreground focus-visible:ring-0"
                disabled={isLoading}
              />
            </div>
          </div>
          <div className="flex justify-between items-center pt-3 border-t border-border">
            <div className="text-sm text-muted-foreground">{content.length > 0 && `${content.length} characters`}</div>
            <Button type="submit" disabled={!content.trim() || isLoading} className="gap-2">
              {isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Send className="h-4 w-4" />}
              Post
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
