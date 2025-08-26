"use server";

import { createScopedLogger } from "@/utils/logger";
import { AuthUser } from "..";
import {  SessionSchema, UserSchema } from "@/types/model";
import { apiWithAuth } from "@/utils/axios";
import z from "zod";
import { ResponseJson } from "@/utils/response";
import { setSessionCookie } from "./cookie";

const logger = createScopedLogger("src:auth:session");

export const validateSession = async (
): Promise<AuthUser | null> => {
  try {
    const schema = z.object({
      session: SessionSchema,
      user: UserSchema,
    })

    const response = await apiWithAuth.get<ResponseJson<z.infer<typeof schema>>>("/api/v1/users/me")
    if (!response.data.success) {
      logger.debug("Validate session failed, err: ", response.data.error);
      return null
    }

    const validated = schema.safeParse(response.data.data)

    if (!validated.success) {
      logger.debug("Validate session response validation failed, err: ", validated.error);
      return null
    }

    const { user, session } = validated.data

    return {
      user,
      session,
    };
  } catch (error) {
    logger.error("Error validating session", error);
    throw new Error("Error validating session");
  }
};