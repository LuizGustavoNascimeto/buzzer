import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createMessageApi } from "../../api/messages";

export function useCreateMessages(groupId, token) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (message) => createMessageApi(message, token),

    onSuccess: (newMessage) => {
      queryClient.setQueryData(["messages", groupId], (old = []) => [
        ...old,
        newMessage,
      ]);
    },
  });
}
