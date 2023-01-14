import request from "supertest";
import {uri} from "../../utils";
import {getOrCreatePool, scheduleDeletion, sql} from "../../db";

const CreateOneToOneConversationMutation = `#graphql
mutation CreateOneToOneConversation($receiverId: Uuid!) {
    createOneToOneConversation(receiverId: $receiverId) {
        id
        draft {
            conversationId
            messageId
            richText {
                text
                parts {
                    __typename
                    ... on PlainRichText {
                        text
                    }
                }
            }
        }
        lastMessage {
            __typename
            ... on TextMessage {
                id
                conversationId
                richText {
                    text
                    parts {
                        __typename
                        ... on PlainRichText {
                            text
                        }
                    }
                }
                updatedAt
                createdAt
            }
        }
        createdAt
    }
}
`;

export const sendCreateOneToOneConversationMutation = async (senderId: string, receiverId: string): Promise<request.Request> => {
    const data = {
        query: CreateOneToOneConversationMutation,
        variables: {receiverId},
    };

    return request(uri)
        .post('/')
        .auth(senderId, {type: 'bearer'})
        .send(data);
}

export const createConversation = async (senderId: string, receiverId: string): Promise<string> => {
    const res = await sendCreateOneToOneConversationMutation(senderId, receiverId);

    expect(res.error).toBeFalsy();

    const conversationId = res.body.data?.createOneToOneConversation.id;
    expect(conversationId).not.toBeNull();

    scheduleDeletion(() => revertCreateConversation(conversationId));

    return conversationId;
}

export const revertCreateConversation = async (conversationId: string): Promise<void> => {
    if (conversationId === undefined) {
        return;
    }

    const pool = await getOrCreatePool();

    await pool.connect(async (connection) => {
        await connection.query(
            sql.typeAlias('void')`
                DELETE
                FROM messenger.conversation
                WHERE id = ${conversationId}
            `
        )
    });
}