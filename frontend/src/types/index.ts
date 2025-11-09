export type MessageCategory = 'EMERGENCY' | 'GENERAL' | 'HELP';

export interface User {
    id: string;
    nickname: string;
    nodeId: string;
}

export interface Message {
    id: string;
    content: string;
    category: MessageCategory;
    sender: User;
    timestamp: string;
}