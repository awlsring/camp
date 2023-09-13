use aws_smithy_http_server::{body::BoxBody, response::{IntoResponse}, proto::{aws_json_10::AwsJson1_0, aws_json_11::AwsJson1_1, aws_json::{runtime_error::RuntimeError}, rest_xml::{RestXml, runtime_error::RuntimeError as XmlError}, rest_json_1::{RestJson1, runtime_error::RuntimeError as RestError}}};

pub struct UnauthorizedException;

impl IntoResponse<AwsJson1_0> for UnauthorizedException {
    fn into_response(self) -> http::Response<BoxBody> {
        IntoResponse::<AwsJson1_0>::into_response(RuntimeError::NotAcceptable)
    }
}

impl IntoResponse<AwsJson1_1> for UnauthorizedException {
    fn into_response(self) -> http::Response<BoxBody> {
        IntoResponse::<AwsJson1_1>::into_response(RuntimeError::NotAcceptable)
    }
}

impl IntoResponse<RestXml> for UnauthorizedException {
    fn into_response(self) -> http::Response<BoxBody> {
        IntoResponse::<RestXml>::into_response(XmlError::NotAcceptable)
    }
}

impl IntoResponse<RestJson1> for UnauthorizedException {
    fn into_response(self) -> http::Response<BoxBody> {
        IntoResponse::<RestJson1>::into_response(RestError::NotAcceptable)
    }
}