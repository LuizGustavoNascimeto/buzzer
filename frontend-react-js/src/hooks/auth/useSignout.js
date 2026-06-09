import { useMutation, useQueryClient } from "@tanstack/react-query";
import { signOut } from "../../api/auth";

export function useSignout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: signOut,
    onSuccess: () => {
      queryClient.resetQueries({ queryKey: ["user"] });
    },
  });
}
