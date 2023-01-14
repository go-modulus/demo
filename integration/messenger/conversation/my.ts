import request from "supertest";
import {uri} from "../../utils";

const MyConversationsQuery = `#graphql
query ListMyConversations {
    myConversations {
        edges {
            cursor
            node {
                __typename
                ... on OneToOneConversation {
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
        }
        pageInfo {
            hasNextPage
            hasPreviousPage
            startCursor
            endCursor
        }
    }
}
`;

export const sendMyConversationsQuery = async (performerId: string): Promise<request.Request> => {
    const data = {
        query: MyConversationsQuery,
        variables: {},
    };

    return request(uri)
        .post('/')
        .auth(performerId, {type: 'bearer'})
        .send(data);
}
