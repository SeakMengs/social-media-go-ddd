"use client"

import type React from "react"

import { useState, useEffect } from "react"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Textarea } from "@/components/ui/textarea"
import { Loader2 } from "lucide-react"
import { Post } from "@/types/model"
import { AuthUser } from "@/auth"
import { updatePost } from "@/api/action"
import { toast } from "sonner"

interface EditPostDialogProps {
  auth: AuthUser
  post: Post | null
  open: boolean
  onOpenChange: (open: boolean) => void
  onPostUpdated?: () => void
}

export function EditPostDialog({ auth, post, open, onOpenChange, onPostUpdated }: EditPostDialogProps) {
  const [content, setContent] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const { user } = auth

  useEffect(() => {
    if (post) {
      setContent(post.content)
    }
  }, [post])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!content.trim() || !post || !user) return

    setIsLoading(true)
    try {
      const response = await updatePost(post.id, {
        content: content.trim(),
      })
      if (response.success) {
        onOpenChange(false)
        onPostUpdated?.()
        toast.success("Post updated successfully")
      } else {
        toast.error(response.error || "Failed to update post")
      }
    } catch (error) {
      toast.error("An unexpected error occurred")
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[525px]">
        <DialogHeader>
          <DialogTitle className="font-heading">Edit Post</DialogTitle>
          <DialogDescription>Make changes to your post. Click save when you're done.</DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit}>
          <div className="py-4">
            <Textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder="What's on your mind?"
              className="min-h-[120px] resize-none"
              disabled={isLoading}
            />
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)} disabled={isLoading}>
              Cancel
            </Button>
            <Button type="submit" disabled={!content.trim() || isLoading}>
              {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Save Changes
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
