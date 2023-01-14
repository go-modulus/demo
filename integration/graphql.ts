export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
    ID: string;
    String: string;
    Boolean: boolean;
    Int: number;
    Float: number;
    MilliTimestamp: any;
    Uuid: any;
    Void: any;
};

export type Conversation = OneToOneConversation;

export type ConversationEdge = {
    __typename?: 'ConversationEdge';
    cursor: Scalars['String'];
    node: Conversation;
};

export type ConversationList = {
    __typename?: 'ConversationList';
    edges: Array<Maybe<ConversationEdge>>;
    pageInfo: PageInfo;
};

export type Draft = {
    __typename?: 'Draft';
    conversationId: Scalars['Uuid'];
    messageId?: Maybe<Scalars['Uuid']>;
    richText?: Maybe<RichText>;
};

export type Message = TextMessage;

export type Mutation = {
    __typename?: 'Mutation';
    createMessage: TextMessage;
    createOneToOneConversation: OneToOneConversation;
    deleteMessage: Scalars['Void'];
    editMessage: TextMessage;
    register: User;
    removeDraft: Scalars['Void'];
    saveDraft: Draft;
};


export type MutationCreateMessageArgs = {
    conversationId: Scalars['Uuid'];
    text?: InputMaybe<Scalars['String']>;
};


export type MutationCreateOneToOneConversationArgs = {
    receiverId: Scalars['Uuid'];
};


export type MutationDeleteMessageArgs = {
    messageId: Scalars['Uuid'];
};


export type MutationEditMessageArgs = {
    messageId: Scalars['Uuid'];
    text?: InputMaybe<Scalars['String']>;
};


export type MutationRegisterArgs = {
    request: RegisterRequest;
};


export type MutationRemoveDraftArgs = {
    conversationId: Scalars['Uuid'];
};


export type MutationSaveDraftArgs = {
    conversationId: Scalars['Uuid'];
    messageId?: InputMaybe<Scalars['Uuid']>;
    text?: InputMaybe<Scalars['String']>;
};

export type OneToOneConversation = {
    __typename?: 'OneToOneConversation';
    createdAt: Scalars['MilliTimestamp'];
    draft?: Maybe<Draft>;
    id: Scalars['Uuid'];
    lastMessage?: Maybe<Message>;
};

export type PageInfo = {
    __typename?: 'PageInfo';
    endCursor: Scalars['String'];
    hasNextPage: Scalars['Boolean'];
    hasPreviousPage: Scalars['Boolean'];
    startCursor: Scalars['String'];
};

export type PlainRichText = {
    __typename?: 'PlainRichText';
    text: Scalars['String'];
};

export type Query = {
    __typename?: 'Query';
    conversation: Conversation;
    myConversations: ConversationList;
    user: User;
    users: UserList;
};


export type QueryConversationArgs = {
    id: Scalars['Uuid'];
};


export type QueryMyConversationsArgs = {
    after?: InputMaybe<Scalars['String']>;
    first?: Scalars['Int'];
};


export type QueryUserArgs = {
    id: Scalars['String'];
};


export type QueryUsersArgs = {
    after?: InputMaybe<Scalars['String']>;
    first: Scalars['Int'];
};

export type RegisterRequest = {
    age: Scalars['Int'];
    email: Scalars['String'];
    name: Scalars['String'];
};

export type RichText = {
    __typename?: 'RichText';
    parts: Array<RichTextPart>;
    text?: Maybe<Scalars['String']>;
};

export type RichTextPart = PlainRichText;

export type TextMessage = {
    __typename?: 'TextMessage';
    conversationId: Scalars['Uuid'];
    createdAt: Scalars['MilliTimestamp'];
    id: Scalars['Uuid'];
    richText?: Maybe<RichText>;
    updatedAt: Scalars['MilliTimestamp'];
};

export type User = {
    __typename?: 'User';
    email: Scalars['String'];
    id: Scalars['String'];
    name: Scalars['String'];
};

export type UserEdge = {
    __typename?: 'UserEdge';
    cursor: Scalars['String'];
    node: User;
};

export type UserList = {
    __typename?: 'UserList';
    edges: Array<UserEdge>;
    pageInfo: PageInfo;
};
