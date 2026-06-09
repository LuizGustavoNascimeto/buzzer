import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createActivityApi } from "../../api/activities";

export function useCreateActivity() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createActivityApi,

    onSuccess: (newActivity) => {
      queryClient.setQueryData(["activities"], (old = []) => [
        newActivity,
        ...old,
      ]);
    },
  });
}
