import { useMutation, useQueryClient } from "@tanstack/react-query";
import { signIn } from "../../api/auth";

export function useSignin() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ email, password }) => signIn(email, password),
    onSuccess: (data) => {
      queryClient.setQueryData(["user"], data); // seta o usuário autenticado
    },
  });
}
