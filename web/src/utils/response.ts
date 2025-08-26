export type SuccessResponse<T> = {
  data: T;
  message: string;
  success: true;
  code: number;
};

export type ErrorResponse<T> = {
  message: string;
  error: string;
  success: false;
  code: number;
};

export type ResponseJson<T = any, E = {}> =
  | SuccessResponse<T>
  | ErrorResponse<E>;

export const responseSuccess = <T = any>(
  message: string,
  data: T = {} as T
): SuccessResponse<T> => ({
  data: data,
  message,
  success: true,
  code: 200,
});

export const responseFailed = (message: string, error: string, code: number) => ({
  message,
  error,
  success: false,
  code,
});

export const responseSomethingWentWrong = (message: string) => {
  return responseFailed(
    message,
    "something went wrong",
    500
  );
};
