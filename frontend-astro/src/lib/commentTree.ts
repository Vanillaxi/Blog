type WithReplies<T> = T & {
  id: number | string;
  nickname: string;
  parentId?: number | string;
  children?: T[];
  replyTo?: {
    id: number | string;
    nickname: string;
  };
};

function idKey(value: number | string | undefined) {
  return String(value ?? "");
}

export function buildTwoLevelCommentTree<T extends WithReplies<T>>(items: T[]): T[] {
  const byID = new Map<string, T>();
  const roots: T[] = [];
  const repliesByRoot = new Map<string, T[]>();

  for (const item of items) {
    byID.set(idKey(item.id), { ...item, children: [...(item.children ?? [])] });
  }

  for (const item of byID.values()) {
    const parentID = idKey(item.parentId);
    if (!parentID || parentID === "0") {
      roots.push(item);
      continue;
    }

    let rootID = parentID;
    let parent = byID.get(parentID);
    while (parent?.parentId && idKey(parent.parentId) !== "0") {
      rootID = idKey(parent.parentId);
      parent = byID.get(rootID);
    }

    const replyTo = byID.get(parentID);
    if (replyTo && !item.replyTo) {
      item.replyTo = { id: replyTo.id, nickname: replyTo.nickname };
    }

    repliesByRoot.set(rootID, [...(repliesByRoot.get(rootID) ?? []), item]);
  }

  return roots.map((root) => ({
    ...root,
    children: [...(root.children ?? []), ...(repliesByRoot.get(idKey(root.id)) ?? [])],
  }));
}
