import axios, {
  AxiosError,
  AxiosResponse,
  HttpStatusCode,
  InternalAxiosRequestConfig,
} from "axios";
import { getApiBaseUrl } from ".";
import { getSessionIdFromCookie } from "@/auth/server/cookie";
import { ResponseJson } from "./response";

const apiBaseUrl = getApiBaseUrl();

export const api = axios.create({
  baseURL: apiBaseUrl,
  // transformResponse: [
  //   (data: string) => {
  //     return data;
  //     if (typeof data !== "string") {
  //       return data;
  //     }

  //     // transform error response to zod formatted error
  //     //   return transformAutoCertErrorToZodFormattedError(JSON.parse(data));
  //   },
  // ],
  validateStatus: (status) => {
    // Resolve only if the status code is less than 500
    // return status < 500;
    return true; // Always resolve the response
  },
});

// Handle auto refresh token and set bearer token in the header. Does not handle auto refresh token for server side call.
export const apiWithAuth = axios.create({
  baseURL: apiBaseUrl,
  // If you want to send cookies with the request. | Keep in mind that api allow origin should not be "*".
  // withCredentials: true,
  // transformResponse: [
  //   (data: string) => {
  //     return data;
  //     //   if (typeof data !== "string") {
  //     //     return data;
  //     //   }

  //     //   // transform error response to zod formatted error
  //     //   return transformAutoCertErrorToZodFormattedError(JSON.parse(data));
  //   },
  // ],
  validateStatus: (status) => {
    // Resolve only if the status code is less than 500
    // return status < 500;
    return true; // Always resolve the response
  },
});

apiWithAuth.interceptors.request.use(
  async (config: InternalAxiosRequestConfig) => {
    const session = await getSessionIdFromCookie();
    if (session) {
      config.headers["Authorization"] = `Bearer ${session}`;
    }
    return config;
  }
);

apiWithAuth.interceptors.response.use(
  (response: AxiosResponse) => {
    return response;
  },
  async (error: AxiosError<ResponseJson<any>>) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & {
      // Extend the _retry such that we can track if the request has already been retried
      _retry?: boolean;
    };
    const status = error.response?.status;
    const isClientSide = typeof window !== "undefined";

    // If the error is 401 and the request has not been retried, try to refresh the token and retry the request
    // It only perform on client side call, because server side call will not have access to the cookie
    if (
      originalRequest &&
      status === HttpStatusCode.Unauthorized &&
      !originalRequest._retry &&
      isClientSide
    ) {
      originalRequest._retry = true;
      //   const {accessToken, refreshToken} = await refreshAccessToken();
      //   if (accessToken && refreshToken) {
      //     setRefreshAndAccessTokenToCookie(refreshToken, accessToken);
      //     return apiWithAuth(originalRequest);
      //   }

      return Promise.reject(error);
    }
    return Promise.reject(error);
  }
);

// const transformAutoCertErrorToZodFormattedError = (data: any) => {
//   if (!data.success) {
//     return {
//       ...data,
//       errors: autocertToFormattedZodError(data.errors as T_AutocertError[]),
//       // message: data.message,
//       // success: false,
//     } satisfies ResponseJson<any, T_ZodErrorFormatted>;
//   }

//   return data;
// };
