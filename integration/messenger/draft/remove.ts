import {uri} from "../../utils";
import request from "supertest";

const RemoveDraftMutation = `#graphql
mutation EditMessage($conversationId: Uuid!) {
    removeDraft(conversationId: $conversationId)
}
`;

export const sendRemoveDraftMutation = async (
    performerId: string,
    conversationId: string,
): Promise<request.Request> => {
    const data = {
        query: RemoveDraftMutation,
        variables: {conversationId},
    };

    return request(uri)
        .post('/')
        .auth(performerId, {type: 'bearer'})
        .send(data);
}

export const removeDraft = async (
    performerId: string,
    conversationId: string,
): Promise<void> => {
    const res = await sendRemoveDraftMutation(performerId, conversationId);

    expect(res.error).toBeFalsy();
}
