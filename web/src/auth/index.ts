import { Session, User } from "@/types/model";
import { IS_PRODUCTION } from "@/utils";

export const SESSION_COOKIE_NAME = "socialmedia";

// Convert session name to __Secure-prefix-type if in production
export function getSessionCookieName(): string {
  return IS_PRODUCTION
    ? `__Secure-${SESSION_COOKIE_NAME}`
    : `${SESSION_COOKIE_NAME}`;
}

export type AuthUser = {
  user: User,
  session: Session
}

export type AuthUserResult = AuthUser | null;