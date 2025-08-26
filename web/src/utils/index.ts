import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export const IS_PRODUCTION = process.env.NODE_ENV === "production";

export const APP_NAME = "social media";

export function getApiBaseUrl(): string {
  return process.env.API_BASE_URL ?? "http://127.0.0.1:8080";
}

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}