import request from "supertest";
import {uri} from "../../utils";

const EditMessageMutation = `#graphql
mutation EditMessage($messageId: Uuid!, $text: String) {
    editMessage(messageId: $messageId, text: $text) {
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

export const sendEditMessage = async (
    performerId: string,
    messageId: string,
    text: string,
): Promise<request.Request> => {
    const data = {
        query: EditMessageMutation,
        variables: {messageId, text},
    };

    return request(uri)
        .post('/')
        .auth(performerId, {type: 'bearer'})
        .send(data);
}

export const editMessage = async (
    performerId: string,
    messageId: string,
    text: string,
): Promise<string> => {
    const res = await sendEditMessage(performerId, messageId, text);

    expect(res.error).toBeFalsy();

    return messageId;
}