import { useQuery } from "@tanstack/react-query";
import { fetchMessages } from "../../api/messages";

export function useMessage(group_id, token) {
  return useQuery({
    queryKey: ["messages", group_id],
    queryFn: () => fetchMessages(group_id, token),
    enabled: !!group_id && !!token,
  });
}