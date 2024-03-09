require "openssl"
require "addressable"
require "excon"

class Mediator
  class << self
    attr_accessor :secret_key
    attr_accessor :base_url

    ENDPOINTS = {
      transform: "/image/transform",
      proxy: "/proxy"
    }

    def url(endpoint, source, file_path, options = {})
      raise "Unsupported endpoint" unless ENDPOINTS.keys.include?(endpoint)

      mediator_path = [ENDPOINTS[endpoint], source, Addressable::URI.escape(file_path)].join("/")
      "#{base_url}#{signed_path(mediator_path, options)}"
    end

    private

    def signed_path(path, options = {})
      uri = Addressable::URI.parse(path)
      uri.query_values = (uri.query_values || {}).merge(options)
      uri.query_values = uri.query_values.merge(s: url_signature(uri.to_s.chomp("?")))
      uri.to_s
    end

    def url_signature(url)
      raise "Missing mediator secret key" if secret_key.blank?

      OpenSSL::HMAC.hexdigest(OpenSSL::Digest.new("sha256"), secret_key, url)
    end
  end
end

# config/initializers/mediator.rb
Rails.application.config.to_prepare do
  Mediator.secret_key = Rails.application.credentials.mediator_secret_key
  Mediator.base_url = Rails.configuration.app.fetch(:mediator_url)
end

# app/helpers/mediator_helper.rb
module MediatorHelper
  def mediator_image_tag(file_path, options = {})
    source = Rails.env.production? ? "images" : "images-dev"
    image_tag Mediator.url(:transform, source, file_path, w: options[:width], h: options[:height]), options
  end
end
