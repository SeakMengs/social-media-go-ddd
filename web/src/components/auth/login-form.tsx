"use client";

import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { toast } from "sonner";
import { Loader2 } from "lucide-react";
import { login } from "@/auth/server/action";
import { createScopedLogger } from "@/utils/logger";
import { wait, WAIT_MS } from "@/utils";

const logger = createScopedLogger("src:components:auth:login-form");

type LoginFormInputs = {
  username: string;
  password: string;
};

export function LoginForm() {
  const router = useRouter();
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    setError,
    clearErrors,
  } = useForm<LoginFormInputs>();

  const onSubmit = async (data: LoginFormInputs) => {
    clearErrors();
    try {
      const response = await login(data);
      if (response.success) {
        toast.success("Login successful!");
        await wait(WAIT_MS);
        router.push("/feed");
      } else {
        setError("root", {
          message:
            response.error || "Login failed. Please check your credentials.",
        });
      }
    } catch {
      setError("root", {
        message: "An unexpected error occurred. Please try again later.",
      });
    }
  };

  return (
    <Card className="w-full max-w-md">
      <CardHeader className="text-center">
        <CardTitle className="text-2xl font-heading">Welcome back</CardTitle>
        <CardDescription>Sign in to your account to continue</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="username">Username</Label>
            <Input
              id="username"
              type="text"
              {...register("username", { required: "Username is required" })}
              disabled={isSubmitting}
            />
            {errors.username && (
              <p className="text-sm text-red-500">{errors.username.message}</p>
            )}
          </div>
          <div className="space-y-2">
            <Label htmlFor="password">Password</Label>
            <Input
              id="password"
              type="password"
              {...register("password", { required: "Password is required" })}
              disabled={isSubmitting}
            />
            {errors.password && (
              <p className="text-sm text-red-500">{errors.password.message}</p>
            )}
          </div>
          {(errors.root?.message ||
            errors.username?.message ||
            errors.password?.message) && (
            <Alert variant="destructive">
              <AlertDescription>
                {errors.root?.message ||
                  errors.username?.message ||
                  errors.password?.message}
              </AlertDescription>
            </Alert>
          )}
          <Button type="submit" className="w-full" disabled={isSubmitting}>
            {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Sign In
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
