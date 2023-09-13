use aws_smithy_http_server::operation::OperationShape;

// use this to validate op to auth when its less painful
pub trait AuthConfig<Op>
where 
    Op: OperationShape,
{
    fn error(&self, err: String) -> <Op as OperationShape>::Error;
}