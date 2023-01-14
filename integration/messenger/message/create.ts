import request from "supertest";
import {uri} from "../../utils";
import {getOrCreatePool, scheduleDeletion, sql} from "../../db";

const CreateMessageMutation = `#graphql
mutation CreateMessage($conversationId: Uuid!, $text: String) {
    createMessage(conversationId: $conversationId, text: $text) {
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
`;

export const sendCreateMessageMutation = async (
    performerId: string,
    conversationId: string,
    text: string,
): Promise<request.Request> => {
    const data = {
        query: CreateMessageMutation,
        variables: {conversationId, text},
    };

    return request(uri)
        .post('/')
        .auth(performerId, {type: 'bearer'})
        .send(data);
}

export const createMessage = async (
    performerId: string,
    conversationId: string,
    text: string,
): Promise<string> => {
    const res = await sendCreateMessageMutation(performerId, conversationId, text);

    expect(res.error).toBeFalsy();

    const messageId = res.body.data?.createMessage.id;
    expect(messageId).not.toBeNull();

    scheduleDeletion(() => revertCreateMessage(messageId))

    return messageId;
}

export const revertCreateMessage = async (messageId: string): Promise<void> => {
    if (messageId === undefined) {
        return;
    }

    const pool = await getOrCreatePool();

    await pool.connect(async (connection) => {
        await connection.query(
            sql.typeAlias('void')`
                DELETE
                FROM messenger.message
                WHERE id = ${messageId}
            `
        )
    });
}