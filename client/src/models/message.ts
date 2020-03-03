export default interface Message {
    id: number;
    userId: number;
    conversationId: number;
    content: string;
    createdAt: Date;
}
