"use client"

import type React from "react"

import { PropsWithChildren, useEffect } from "react"
import { useRouter, usePathname } from "next/navigation"
import { Navigation } from "./navigation"
import { AuthUserResult } from "@/auth"
import { Toaster } from "sonner"

interface AppLayoutProps {
  auth: AuthUserResult
}

export function AppLayout({ children, auth }: PropsWithChildren<AppLayoutProps>) {
  const isAuthenticated = auth && auth.user ? true : false
  const router = useRouter()
  const pathname = usePathname()

  // const publicRoutes = ["/", "/login", "/register"]
  // const isPublicRoute = publicRoutes.includes(pathname)

  // useEffect(() => {
  //   if (!isAuthenticated && !isPublicRoute) {
  //     router.push("/login")
  //   }
  // }, [isAuthenticated, isPublicRoute, router])

  // if (!isAuthenticated && !isPublicRoute) {
  //   return null
  // }

  return (
    <div className="min-h-screen bg-background">
      {auth && <Navigation auth={auth} />}
      <main className={isAuthenticated ? "" : ""}>{children}</main>
      <Toaster />
    </div>
  )
}
