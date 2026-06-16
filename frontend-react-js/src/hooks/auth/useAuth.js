import { useQuery } from "@tanstack/react-query";
import { checkAuth } from "../../api/auth";

export function useAuth() {
  return useQuery({
    queryKey: ["user"],
    queryFn: () =>  checkAuth(),
    select: (data) => ({
      display_name: data?.attributes?.preferred_username,
      handle: data?.attributes?.name,
      token: data?.token,
    }),
  });
}
