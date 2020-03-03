import Message from '@/models/message';

export default interface Conversation {
    id: number;
    createdAt: Date;
    updatedAt: Date;
    userIds: number[];
    messages: Message[];
    firstLoad: boolean;
    firstMessage: boolean;
}
