"use server";

import { createScopedLogger } from "@/utils/logger";
import { AuthUserResult } from "..";
import { validateSession } from "./session";
import { cache } from "react";
import { api, apiWithAuth } from "@/utils/axios";
import { ResponseJson, responseSomethingWentWrong } from "@/utils/response";
import { z } from "zod";
import { setSessionCookie } from "./cookie";
import { SessionSchema, UserSchema } from "@/types/model";

const logger = createScopedLogger("src:auth:action");

export const getAuthUser = cache(async (): Promise<AuthUserResult> => {
  try {
    const result = await validateSession();
    if (!result) {
      logger.debug("Session is invalid or expired");
      return null;
    }

    return result;
  } catch (error) {
    logger.error("Error getting auth user", error);
  }

  return null;
});

const loginResponseSchema = z.object({
  session: SessionSchema,
  user: UserSchema,
});
export async function login(data: {
  username: string;
  password: string;
}): Promise<ResponseJson<z.infer<typeof loginResponseSchema>>> {
  try {
    const response = await api.post<
      ResponseJson<z.infer<typeof loginResponseSchema>>
    >("/api/v1/auth/login", data);
    if (!response.data.success) {
      logger.debug("Login failed, err: ", response.data.error);
      return response.data;
    }

    const validated = loginResponseSchema.safeParse(response.data.data);
    if (!validated.success) {
      logger.debug("Login response validation failed, err: ", validated.error);
      return response.data;
    }

    const { session } = validated.data;

    await setSessionCookie(session.id, new Date(session.expireAt));

    return response.data;
  } catch (error) {
    logger.error("Error logging in", error);
    return responseSomethingWentWrong("Error logging in (catch)");
  }
}

export async function logout(): Promise<void> {
  try {
    const response = await apiWithAuth.delete("/api/v1/auth/logout");
    if (!response.data.success) {
      logger.debug("Logout failed", response.data.error);
      return;
    }

    logger.debug("User logged out successfully");
  } catch (error) {
    logger.error("Error logging out", error);
  }
}
