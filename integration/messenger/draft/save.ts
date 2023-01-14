import request from "supertest";
import {uri} from "../../utils";
import {getOrCreatePool, scheduleDeletion, sql} from "../../db";

const SaveDraftMutation = `#graphql
mutation SaveDraft($conversationId: Uuid!, $messageId: Uuid, $text: String) {
    saveDraft(conversationId: $conversationId, messageId: $messageId, text: $text) {
        messageId
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
    }
}
`;

export const sendSaveDraftMutation = async (
    performerId: string,
    conversationId: string,
    messageId: string | null,
    text: string,
): Promise<request.Request> => {
    const data = {
        query: SaveDraftMutation,
        variables: {conversationId, messageId, text},
    };

    return request(uri)
        .post('/')
        .auth(performerId, {type: 'bearer'})
        .send(data);
}

export const saveDraft = async (
    performerId: string,
    conversationId: string,
    messageId: string | null,
    text: string,
): Promise<void> => {
    const res = await sendSaveDraftMutation(
        performerId, conversationId,
        messageId,
        text
    );

    expect(res.error).toBeFalsy();

    scheduleDeletion(() => revertSaveDraft(performerId, conversationId));
}

export const revertSaveDraft = async (authorId: string, conversationId: string): Promise<void> => {
    if (authorId === undefined || conversationId === undefined) {
        return;
    }

    const pool = await getOrCreatePool();

    await pool.connect(async (connection) => {
        await connection.query(
            sql.typeAlias('void')`
                DELETE
                FROM messenger.draft
                WHERE author_id = ${authorId} AND conversation_id = ${conversationId}
            `
        )
    });
}