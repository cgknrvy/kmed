import { useAuthStore } from "#/stores/auth.ts";

const API_BASE_URL = "/api/v1/";

let refreshPromise: Promise<string> | null = null;

/**
 * Makes a fetch request to the given url with the options passed.
 * The `Authorization` header is automatically added to ensure that the requests are
 * valid.
 *
 * If the first request returns an authorization error(401), the access token is refreshed
 * and then the request is resent with the new access token. If getting a new access token
 * is impossible then an error is returned. If an error is returned, it is advised to redirect
 * the user to the login page to have them login and get new tokens.
 *
 * @param url url to send the request to
 * @param options set of options passed to the fetch API
 */
export async function apiFetchWithRefresh(
  url: string | URL,
  options: RequestInit = {},
) {
  const accessToken = useAuthStore.getState().accessToken;
  if (typeof url === "string" && !url.startsWith(API_BASE_URL)) {
    url = joinApiBase(url);
  }

  const response = await fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      Authorization: `Bearer ${accessToken}`,
    },
  });

  if (response.status !== 401) {
    return response;
  }

  if (!refreshPromise) {
    refreshPromise = refreshAccessToken().finally(
      // reset refreshPromise after the request returns even if it fails
      () => (refreshPromise = null),
    );
  }

  const token = await refreshPromise;

  return fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      Authorization: `Bearer ${token}`,
    },
  });
}

/**
 * Makes a fetch request to the given url with the options passed.
 * The `Authorization` header is automatically added to ensure that the requests are
 * valid.
 *
 * Does not refresh the access token if its invalid.
 *
 * @param url url to send the request to
 * @param options set of options passed to the fetch API
 */
export async function apiFetchWithoutRefresh(
  url: string | URL,
  options: RequestInit = {},
) {
  const accessToken = useAuthStore.getState().accessToken;
  if (typeof url === "string" && !url.startsWith(API_BASE_URL)) {
    url = joinApiBase(url);
  }

  return fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      Authorization: `Bearer ${accessToken}`,
    },
  });
}

/**
 * Sends a refresh request to the api and gets a new access token if the current
 * refresh token is valid. Sets the retrieved access token in the `useAuthStore`
 * accessToken state.
 */
export async function refreshAccessToken() {
  const url = joinApiBase("auth/refresh");
  const response = await fetch(url, {
    method: "POST",
    credentials: "include",
  });

  if (!response.ok) {
    throw new ApiError(response.status, "unable to refresh access token");
  }
  const { accessToken } = await response.json();

  useAuthStore.getState().setAccessToken(accessToken);

  return accessToken;
}

export async function apiFetchUnauthorized(
  url: string | URL,
  options: RequestInit = {},
) {
  if (typeof url === "string" && !url.startsWith(API_BASE_URL)) {
    url = joinApiBase(url);
  }
  return fetch(url, {
    ...options,
    headers: {
      ...options.headers,
    },
  });
}

function joinApiBase(path: string | URL) {
  const url = URL.parse(path, API_BASE_URL);
  if (url === null) {
    return API_BASE_URL + path;
  }
  return url;
}

export class ApiError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.status = status;
    this.name = "ApiError";
  }
}
